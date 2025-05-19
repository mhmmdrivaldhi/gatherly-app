package models

import "time"

type Event struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   		string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsPaid      bool      `json:"is_paid"`
	Tickets     []Ticket  `gorm:"foreignKey:EventID"`
	Capacity    int       `json:"capacity"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	PosterURL   string    `json:"poster_url"`
	Status      string    `json:"status"`
}
