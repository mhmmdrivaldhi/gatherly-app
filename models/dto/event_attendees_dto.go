package dto

type AttendeeRegisterRequest struct {
    EventID      int    `json:"eventId" binding:"required"`
    TicketTypeID int    `json:"ticketTypeId" binding:"required"` // <-- Make sure this line exists
    RSVPStatus   string `json:"rsvpStatus" binding:"required,oneof=pending attending not_attending maybe"`
}

type AttendeeCancelRequest struct {
	UserID  int `json:"userId" binding:"required"`
	EventID int `json:"eventId" binding:"required"`
}

type AttendeePaymentRequest struct {
	UserID  int `json:"userId" binding:"required"`
	EventID int `json:"eventId" binding:"required"`
}

type AttendeeRSVPRequest struct {
	UserID    int    `json:"userId" binding:"required"`
	EventID   int    `json:"eventId" binding:"required"`
	NewStatus string `json:"newStatus" binding:"required,oneof=going interested not_going"`
}
