package dto

import "time"

type CreateEventRequestDTO struct {
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	IsPaid      bool    `json:"is_paid"`
	Price       float64 `json:"price"`
	Capacity    int     `json:"capacity" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	PosterURL   string  `json:"poster_url"`
	Status      string  `json:"status"`
}

type UpdateEventRequestDTO struct {
	Name        *string  `json:"name"`
	Category    *string  `json:"category"`
	Description *string  `json:"description"`
	StartDate   *string  `json:"start_date"`
	EndDate     *string  `json:"end_date"`
	IsPaid      *bool    `json:"is_paid"`
	Price       *float64 `json:"price"`
	Capacity    *int     `json:"capacity"`
	Address     *string  `json:"address"`
	PosterURL   *string  `json:"poster_url"`
	Status      *string  `json:"status"`
}

type EventResponseDTO struct {
	Name        string             `json:"name"`
	Category    string             `json:"category"`
	Description string             `json:"description"`
	StartDate   string             `json:"start_date"`
	EndDate     string             `json:"end_date"`
	IsPaid      bool               `json:"is_paid"`
	Ticket      *TicketResponseDTO `json:"ticket"`
	Capacity    int                `json:"capacity"`
	Latitude    float64            `json:"latitude"`
	Longitude   float64            `json:"longitude"`
	PosterURL   string             `json:"poster_url"`
	Status      string             `json:"status"`
}

type EventNearbyDistanceResponseDTO struct {
	ID 			int		  `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time    `json:"end_date"`
	IsPaid      bool      `json:"is_paid"`
	Ticket      *TicketResponseDTO `json:"ticket"`
	Capacity    int       `json:"capacity"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Distance    float32   `json:"distance"`
	PosterURL   string    `json:"poster_url"`
	Status      string    `json:"status"`
}

type GeneralResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}