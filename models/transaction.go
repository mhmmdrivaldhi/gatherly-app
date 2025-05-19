package models

import (
	"time"

	"gorm.io/gorm"
)

type Transactions struct {
	gorm.Model
	UserId                      int       `json:"user_id" gorm:"not null;index"`
	EventId                     int       `json:"event_id" gorm:"not null;index"`
	Event                       Event    `gorm:"foreignKey:EventId;references:ID"`
	Amount                      float64   `json:"amount" gorm:"not null"`
	TransactionDate             time.Time `json:"transaction_date" gorm:"not null"`
	Status                      string    `json:"status" gorm:"not null"`
	PaymentMethod               string    `json:"payment_method" gorm:"not null"`
	PaymentGatewayTransactionId string    `json:"payment_gateway_transaction_id" gorm:"not null"`
	Items                       string    `json:"items" gorm:"not null"`
	Notes                       string    `json:"notes"`
	Url                         string    `json:"url"`
}

