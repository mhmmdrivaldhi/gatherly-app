package models

import (
	"time"
)

type EventAttendee struct {
	UserID        int        `json:"user_id"`
	EventID       int        `json:"event_id"`
	TicketTypeID  *int       `json:"ticket_type_id"`     // <-- ADD THIS LINE
	Event         Event     `gorm:"foreignKey:EventID"` // <-- Tambahkan relasi ke Events
	RSVPStatus    string     `json:"rsvp_status"`
	RSVPDate      *time.Time `json:"rsvp_date,omitempty"`
	PaymentStatus string     `json:"payment_status"`
	TicketCode    *string    `json:"ticket_code,omitempty"`
}
