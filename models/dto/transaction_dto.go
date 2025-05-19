package dto

import "time"

type CreateTransaction struct {
	UserId          int        `json:"user_id" binding:"required"`
	EventId         int        `json:"event_id" binding:"required"`
	TransactionDate *time.Time `json:"transaction_date" gorm:"not null"`
	Amount          float64    `json:"amount" binding:"required"`
	Items           string     `json:"items" binding:"required"`
	Notes           string     `json:"notes"`
}

type GetTransactionsByDate struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type GetTransactionsByAmount struct {
	MinAmount float64 `json:"min_amount"`
	MaxAmount float64 `json:"max_amount"`
}

type MidtransSnapReq struct {
	TransactionDetails struct {
		OrderID     string `json:"order_id"`
		GrossAmount int    `json:"gross_amount"`
	} `json:"transaction_details"`
	Customer string `json:"customer"`
	Items    string `json:"items_details"`
}

type MidtransSnapResp struct {
	Token        string `json:"token"`
	RedirectURL  string `json:"redirect_url"`
	ErrorMessage string `json:"error_message"`
}

type MidtransNotification struct {
	StatusCode        string `json:"status_code"`
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
}
