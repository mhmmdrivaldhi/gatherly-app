package repositories

import (
	"errors" // Import errors package
	"fmt"
	"gatherly-app/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause" // Import for locking feature
)

// --- Interface Definition ---
type TicketRepository interface {
	// Existing methods
	Save(ticket []models.Ticket) ([]models.Ticket, error)
	Delete(idTicket []int) (models.Ticket, error) // Note: Returns empty Ticket on success
	FindById(eventId int) ([]models.Ticket, error) // Assumes this finds tickets by EVENT ID

	// --- ADDED METHODS ---
	FindTicketByID(id int) (*models.Ticket, error)          // Find a single ticket type by its primary key ID
	FindTicketByIDForUpdate(id int) (*models.Ticket, error) // Find a single ticket type by ID and lock the row
	DecrementQuota(id int) error                           // Decrease quota for a specific ticket ID
	// ---------------------
}

// --- Struct Definition ---
type ticketRepositoryImpl struct {
	db *gorm.DB
}

// --- Constructor ---
// Ensure the constructor returns the interface type
func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepositoryImpl{db: db}
}

// --- Existing Method Implementations ---

func (t *ticketRepositoryImpl) Save(ticket []models.Ticket) ([]models.Ticket, error) {
	if len(ticket) == 0 {
		// It's okay to try saving zero tickets, GORM handles this. Return empty slice.
		return []models.Ticket{}, nil
		// return []models.Ticket{}, fmt.Errorf("no tickets to save") // Or return error if preferred
	}
	// Use Create for new records ensure primary keys aren't specified beforehand if auto-incrementing
	res := t.db.Create(&ticket)

	if res.Error != nil {
		return ticket, fmt.Errorf("error creating ticket types: %w", res.Error)
	}

	return ticket, nil
}

func (t *ticketRepositoryImpl) Delete(idTicket []int) (models.Ticket, error) {
	var ticket models.Ticket // Placeholder variable needed for GORM Delete signature
	if len(idTicket) == 0 {
		return ticket, fmt.Errorf("no IDs provided for deletion")
	}

	// Delete records from models.Ticket table where 'id' is in the slice idTicket
	res := t.db.Where("id IN ?", idTicket).Delete(&models.Ticket{})
	if res.Error != nil {
		return ticket, fmt.Errorf("failed to delete ticket types: %w", res.Error)
	}

	// Optional: Log if rows affected is 0, but don't return error unless deletion failure is critical
	if res.RowsAffected == 0 {
		fmt.Printf("Warning: No ticket types deleted for IDs: %v (they might not have existed)\n", idTicket)
	}

	// Return empty ticket struct on success, as Delete doesn't retrieve deleted records
	return ticket, nil
}

// FindById finds all ticket types associated with a specific EVENT ID
func (t *ticketRepositoryImpl) FindById(eventId int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	// Find all tickets where event_id matches
	res := t.db.Where("event_id = ?", eventId).Find(&tickets)
	if res.Error != nil {
		// Log the error appropriately in a real application
		fmt.Printf("Error finding tickets for event ID %d: %v\n", eventId, res.Error)
		return nil, res.Error // Return the error
	}
	// Return the slice (will be empty if no tickets found, which is not an error)
	return tickets, nil
}

// --- IMPLEMENTATION OF ADDED METHODS ---

// FindTicketByID finds a single ticket type by its primary key ID
func (t *ticketRepositoryImpl) FindTicketByID(id int) (*models.Ticket, error) {
	var ticket models.Ticket
	result := t.db.First(&ticket, id) // Find by primary key 'id'
	if result.Error != nil {
		// Return gorm.ErrRecordNotFound directly if that's the error
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		// For other errors, wrap them
		return nil, fmt.Errorf("error finding ticket by ID %d: %w", id, result.Error)
	}
	return &ticket, nil
}

// FindTicketByIDForUpdate finds a single ticket type by ID and locks the row for update
func (t *ticketRepositoryImpl) FindTicketByIDForUpdate(id int) (*models.Ticket, error) {
	var ticket models.Ticket
	// Use Clauses(clause.Locking{Strength: "UPDATE"}) for pessimistic locking
	// This ensures the row is locked until the current transaction completes
	result := t.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&ticket, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error // Return specific error for not found
		}
		return nil, fmt.Errorf("error finding and locking ticket by ID %d: %w", id, result.Error)
	}
	return &ticket, nil
}

// DecrementQuota decreases the quota for a specific ticket ID by 1 atomically
func (t *ticketRepositoryImpl) DecrementQuota(id int) error {
	// Use UpdateColumn with gorm.Expr for atomic update `quota = quota - 1`
	// Add `quota > 0` condition to prevent decrementing below zero
	result := t.db.Model(&models.Ticket{}).Where("id = ? AND quota > 0", id).UpdateColumn("quota", gorm.Expr("quota - ?", 1))

	if result.Error != nil {
		return fmt.Errorf("error decrementing quota for ticket ID %d: %w", id, result.Error)
	}

	// Check if any row was actually updated
	if result.RowsAffected == 0 {
		// If no rows affected, check if the ticket exists to determine if quota was already 0
		_, findErr := t.FindTicketByID(id) // Use the method we just defined
		if findErr != nil {
			// If FindTicketByID returned RecordNotFound, the ticket doesn't exist
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return fmt.Errorf("cannot decrement quota: ticket ID %d not found", id)
			}
			// Some other error occurred during the check
			return fmt.Errorf("error checking ticket existence after failed quota decrement (ID %d): %w", id, findErr)
		}
		// If the ticket exists but no rows were affected, quota must have been 0
		return fmt.Errorf("cannot decrement quota: ticket ID %d already has zero quota", id)
	}

	// Quota decremented successfully
	return nil
}

// --- END OF ADDED METHODS ---