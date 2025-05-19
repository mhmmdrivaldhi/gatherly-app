package dto

type PayloadTicket struct {
	Ids []int `json:"ids" binding:"required"`
}

type TicketResponseDTO struct {
	ID int `json:"id"`
	TicketType string `json:"ticketType"`
	Price int `json:"price"`
	Quota int `json:"quota"`
	Status string `json:"status"`
}
