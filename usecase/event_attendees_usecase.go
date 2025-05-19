package usecase

import (
	"context"
	"errors"
	"fmt"
	"gatherly-app/models"
	"gatherly-app/models/dto"
	"gatherly-app/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Interface Definition ---
// Added ticketTypeID to Register signature
type EventAttendeeUseCase interface {
	Register(ctx context.Context, userID, eventID, ticketTypeID int, rsvpStatus string) (*models.EventAttendee, error)
	CancelRegistration(ctx context.Context, userID, eventID int) error
	GetRegistrationDetails(ctx context.Context, userID, eventID int) (*models.EventAttendee, error)
	ListAttendeesForEvent(ctx context.Context, eventID int) ([]*models.EventAttendee, error)
	ListUserRegistrations(ctx context.Context, userID int) ([]*models.EventAttendee, error)
	ConfirmPayment(ctx context.Context, userID, eventID int) (*models.EventAttendee, error)
	UpdateRSVPStatus(ctx context.Context, userID, eventID int, newStatus string) (*models.EventAttendee, error)
}

// --- Struct Definition ---
// Added new repositories and use case dependencies
type eventAttendeeUseCaseImpl struct {
	attendeeRepo  repositories.EventAttendeeRepository
	eventRepo     repositories.EventsRepository // Added
	ticketRepo    repositories.TicketRepository // Added
	transactionUC TransactionUsecase          // Added
}

// --- Constructor ---
// Updated constructor to accept new dependencies
func NewEventAttendeeUseCase(
	attendeeRepo repositories.EventAttendeeRepository,
	eventRepo repositories.EventsRepository, // Added
	ticketRepo repositories.TicketRepository, // Added
	transactionUC TransactionUsecase, // Added
) EventAttendeeUseCase {
	return &eventAttendeeUseCaseImpl{
		attendeeRepo:  attendeeRepo,
		eventRepo:     eventRepo,     // Initialized
		ticketRepo:    ticketRepo,    // Initialized
		transactionUC: transactionUC,
	}
}

// --- Helper Function ---
func generateTicketCode(eventID, userID int) string {
	return uuid.NewString()
}

// --- Register Method (Modified) ---
// Updated Register method signature and logic
func (uc *eventAttendeeUseCaseImpl) Register(ctx context.Context, userID, eventID, ticketTypeID int, rsvpStatus string) (*models.EventAttendee, error) {

	// --- Basic Input Validation ---
	allowedRSVP := map[string]bool{"pending": true, "attending": true, "not_attending": true, "maybe": true}
	if !allowedRSVP[rsvpStatus] {
		return nil, fmt.Errorf("invalid RSVP status: %s", rsvpStatus)
	}

	// --- Check Existing Registration ---
	// Note: This simple check might need refinement if users can change ticket types.
	existingAttendee, err := uc.attendeeRepo.FindByUserAndEvent(ctx, userID, eventID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking existing registration: %w", err)
	}
	if existingAttendee != nil {
		// Handle update scenario if needed, or return error if re-registration is disallowed.
		// For now, let's prevent re-registration for simplicity.
		return nil, fmt.Errorf("user %d is already registered for event %d", userID, eventID)
		/* // Example: Update existing (more complex logic might be needed)
		existingAttendee.RSVPStatus = rsvpStatus
		existingAttendee.RSVPDate = &now
		// Potentially update TicketTypeID? Needs careful consideration.
		err = uc.attendeeRepo.Update(ctx, existingAttendee)
		// ... error handling ...
		return existingAttendee, nil
		*/
	}

	// --- Step 1: Fetch and Validate Ticket Type ---
	// IMPORTANT: Assumes ticketRepo has or will have a method FindTicketByID(id int) (*models.Ticket, error)
	// Adjust if your method signature is different (e.g., if FindById returns []models.Ticket)
	ticketType, err := uc.ticketRepo.FindTicketByID(ticketTypeID) // <<--- ADJUST THIS METHOD CALL IF NEEDED
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ticket type with ID %d not found", ticketTypeID)
		}
		return nil, fmt.Errorf("error fetching ticket type details: %w", err)
	}
	if ticketType == nil { // Defensive check
		return nil, fmt.Errorf("ticket type with ID %d not found (nil returned)", ticketTypeID)
	}

	// Validate ticket belongs to the correct event
	if ticketType.EventID != eventID {
		return nil, fmt.Errorf("ticket type ID %d does not belong to event ID %d", ticketTypeID, eventID)
	}
	// Validate quota
	if ticketType.Quota <= 0 {
		return nil, fmt.Errorf("ticket type '%s' is sold out", ticketType.TicketType)
	}
	// Validate status (adjust "available" if you use different status strings)
	if ticketType.Status != "available" {
		return nil, fmt.Errorf("ticket type '%s' is not currently available for purchase", ticketType.TicketType)
	}

	// --- Step 2: Fetch Event Details ---
	event, err := uc.eventRepo.FindEventByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event with ID %d not found", eventID)
		}
		return nil, fmt.Errorf("error fetching event details: %w", err)
	}

	// --- Step 3: Determine Initial Payment Status ---
	paymentStatus := "unpaid" // Default for free events
	if event.IsPaid {
		paymentStatus = "pending" // Requires payment
	}

	// --- Step 4: Create EventAttendee Record ---
	// CONSIDER WRAPPING Steps 4 & 5 IN A DATABASE TRANSACTION FOR ATOMICITY
	now := time.Now()
	newAttendee := &models.EventAttendee{
		UserID:        userID,
		EventID:       eventID,
		TicketTypeID:  &ticketTypeID, // Link to the specific ticket type
		RSVPStatus:    rsvpStatus,
		RSVPDate:      &now,
		PaymentStatus: paymentStatus,
		TicketCode:    nil, // Ticket code generated upon successful payment confirmation
	}

	err = uc.attendeeRepo.Create(ctx, newAttendee)
	if err != nil {
		// If using DB transaction, rollback here
		return nil, fmt.Errorf("failed to create registration record: %w", err)
	}

	// --- Step 5: Create Transaction if Event is Paid ---
	if event.IsPaid {
		transactionInput := dto.CreateTransaction{
			UserId:          userID,
			EventId:         eventID,
			TransactionDate: &now,
			Amount:          float64(ticketType.Price),                                         // Use price from the specific ticket type
			Items:           fmt.Sprintf("Ticket: %s (%s)", event.Name, ticketType.TicketType), // Descriptive item name
			Notes:           fmt.Sprintf("Auto-created for registration EventID: %d, TicketTypeID: %d", eventID, ticketTypeID),
		}

		// Prepare Midtrans request details (OrderID needs to be unique)
		midtransOrderID := uuid.NewString() // Generate a unique ID for Midtrans
		midtransInput := dto.MidtransSnapReq{
			TransactionDetails: struct {
				OrderID     string `json:"order_id"`
				GrossAmount int    `json:"gross_amount"`
			}{
				OrderID:     midtransOrderID,
				GrossAmount: ticketType.Price, // Use the integer price for Midtrans GrossAmount
			},
			// Add Customer and Item details as required by your MidtransService and DTO
			Customer: fmt.Sprintf("User ID: %d", userID), // Example customer detail
			Items:    transactionInput.Items,             // Can reuse the items string
		}

		// Call the injected Transaction Use Case
		_, txErr := uc.transactionUC.CreateTransaction(transactionInput, midtransInput) // Pass both DTOs
		if txErr != nil {
			// CRITICAL: Registration saved, payment failed. Requires robust handling.
			// If using DB transaction, rollback attendeeRepo.Create before returning error.
			// Log the detailed error for debugging.
			fmt.Printf("ERROR: Registration for UserID %d, EventID %d saved, but failed to initiate transaction: %v\n", userID, eventID, txErr)
			// Return the attendee record but indicate the payment initiation failure clearly.
			// The user might need to retry payment later or contact support.
			// Returning a specific error might be better for the controller to handle.
			return newAttendee, fmt.Errorf("registration successful, but failed to start payment process: %w. Please try paying later or contact support", txErr)
		}
		// Transaction initiation successful (Midtrans link/token generated by transactionUC)
		// The transaction record itself is created within transactionUC.CreateTransaction
	}

	// If using DB transaction, commit here

	return newAttendee, nil // Registration successful (payment initiated if applicable)
}

// --- CancelRegistration Method (Unchanged) ---
func (uc *eventAttendeeUseCaseImpl) CancelRegistration(ctx context.Context, userID, eventID int) error {
	_, err := uc.attendeeRepo.FindByUserAndEvent(ctx, userID, eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("registration not found for user %d, event %d", userID, eventID)
		}
		return fmt.Errorf("error checking registration before delete: %w", err)
	}

	// TODO: Add logic here to check if a transaction exists for this registration
	// and potentially cancel it via transactionUC if it's still pending.

	err = uc.attendeeRepo.Delete(ctx, userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete registration: %w", err)
	}

	return nil
}

// --- GetRegistrationDetails Method (Unchanged) ---
func (uc *eventAttendeeUseCaseImpl) GetRegistrationDetails(ctx context.Context, userID, eventID int) (*models.EventAttendee, error) {
	attendee, err := uc.attendeeRepo.FindByUserAndEvent(ctx, userID, eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("registration not found")
		}
		return nil, fmt.Errorf("failed to retrieve registration details: %w", err)
	}
	if attendee == nil {
		return nil, errors.New("registration not found")
	}
	return attendee, nil
}

// --- ListAttendeesForEvent Method (Unchanged) ---
func (uc *eventAttendeeUseCaseImpl) ListAttendeesForEvent(ctx context.Context, eventID int) ([]*models.EventAttendee, error) {
	attendees, err := uc.attendeeRepo.ListByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attendees for event %d: %w", eventID, err)
	}
	return attendees, nil
}

// --- ListUserRegistrations Method (Unchanged) ---
func (uc *eventAttendeeUseCaseImpl) ListUserRegistrations(ctx context.Context, userID int) ([]*models.EventAttendee, error) {
	attendees, err := uc.attendeeRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list registrations for user %d: %w", userID, err)
	}
	return attendees, nil
}

// --- ConfirmPayment Method ---
// Needs modification to handle quota decrement
func (uc *eventAttendeeUseCaseImpl) ConfirmPayment(ctx context.Context, userID, eventID int) (*models.EventAttendee, error) {
	// --- Fetch Attendee Record ---
	attendee, err := uc.attendeeRepo.FindByUserAndEvent(ctx, userID, eventID)
	if err != nil {
		// Handle not found and other errors...
		return nil, fmt.Errorf("failed to find registration for payment confirmation: %w", err)
	}
	if attendee == nil {
		return nil, errors.New("cannot confirm payment for non-existent registration")
	}

	// --- Check if already paid ---
	if attendee.PaymentStatus == "paid" && attendee.TicketCode != nil {
		return attendee, nil // Already confirmed
	}

	// --- Check if TicketTypeID is linked ---
	if attendee.TicketTypeID == nil {
		return nil, errors.New("cannot confirm payment: registration is not linked to a specific ticket type")
	}

	// --- Decrement Quota (CRITICAL: Use DB Transaction) ---
	// It's highly recommended to wrap this section in a database transaction
	// Start transaction here...

	// Fetch the specific ticket type with locking to prevent race conditions
	ticketType, err := uc.ticketRepo.FindTicketByIDForUpdate(*attendee.TicketTypeID) // <<--- NEEDS A REPO METHOD WITH LOCKING (e.g., SELECT ... FOR UPDATE)
	if err != nil {
		// Rollback transaction if started
		return nil, fmt.Errorf("failed to lock ticket type for quota update: %w", err)
	}
	if ticketType == nil {
		// Rollback transaction if started
		return nil, fmt.Errorf("ticket type %d not found during payment confirmation", *attendee.TicketTypeID)
	}

	// Double-check quota before decrementing
	if ticketType.Quota <= 0 {
		// Rollback transaction if started
		// Payment might have succeeded via Midtrans, but quota ran out.
		// This indicates a potential issue (overselling). Log this situation.
		// Update attendee status to maybe "payment_received_no_quota"?
		attendee.PaymentStatus = "failed_no_quota" // Example status
		uc.attendeeRepo.Update(ctx, attendee)      // Update attendee status
		// Commit transaction here (to save the 'failed_no_quota' status)
		return nil, fmt.Errorf("payment confirmed, but ticket type '%s' is now sold out", ticketType.TicketType)
	}

	// Decrement the quota
	err = uc.ticketRepo.DecrementQuota(*attendee.TicketTypeID) // <<--- NEEDS A REPO METHOD TO UPDATE quota = quota - 1
	if err != nil {
		// Rollback transaction if started
		return nil, fmt.Errorf("failed to decrement ticket quota: %w", err)
	}

	// --- Update Attendee Record ---
	attendee.PaymentStatus = "paid"
	attendee.TicketCode = new(string)
	*attendee.TicketCode = generateTicketCode(eventID, userID)
	now := time.Now()
	attendee.RSVPDate = &now // Update RSVP date to reflect payment confirmation time

	err = uc.attendeeRepo.Update(ctx, attendee)
	if err != nil {
		// Rollback transaction if started
		return nil, fmt.Errorf("failed to update registration after payment confirmation: %w", err)
	}

	// Commit transaction here...

	return attendee, nil
}

// --- UpdateRSVPStatus Method (Unchanged for now) ---
func (uc *eventAttendeeUseCaseImpl) UpdateRSVPStatus(ctx context.Context, userID, eventID int, newStatus string) (*models.EventAttendee, error) {
	allowedRSVP := map[string]bool{"pending": true, "attending": true, "not_attending": true, "maybe": true, "going": true, "interested": true, "not_going": true} // Expanded allowed statuses
	if !allowedRSVP[newStatus] {
		return nil, fmt.Errorf("invalid RSVP status: %s", newStatus)
	}

	attendee, err := uc.attendeeRepo.FindByUserAndEvent(ctx, userID, eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cannot update status for non-existent registration")
		}
		return nil, fmt.Errorf("failed to find registration for status update: %w", err)
	}
	if attendee == nil {
		return nil, errors.New("cannot update status for non-existent registration")
	}

	attendee.RSVPStatus = newStatus
	now := time.Now()
	attendee.RSVPDate = &now

	err = uc.attendeeRepo.Update(ctx, attendee)
	if err != nil {
		return nil, fmt.Errorf("failed to update RSVP status: %w", err)
	}

	return attendee, nil
}
