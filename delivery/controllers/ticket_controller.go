package controllers

import (
	"fmt"
	"gatherly-app/models"
	"gatherly-app/models/dto"
	"gatherly-app/usecase"
	"gatherly-app/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketUseCase usecase.TicketUseCase
	rg            *gin.RouterGroup
}

func NewTicketController(ticketUseCase usecase.TicketUseCase, rg *gin.RouterGroup) *TicketController {
	return &TicketController{ticketUseCase: ticketUseCase, rg: rg}
}

func (tc *TicketController) Route() {
	tc.rg.POST("/ticket", tc.createTicket)
	tc.rg.DELETE("/ticket", tc.deleteTicket)
}

// @Summary Create tickets
// @Description Creates a batch of tickets for an event
// @Tags tickets
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param tickets body []model.Ticket true "List of tickets to create"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/ticket [post]
// @Security BearerAuth
func (tc *TicketController) createTicket(ctx *gin.Context) {
	var payload []models.Ticket

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	ticket, err := tc.ticketUseCase.CreateTicket(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ticketResponse := make([]gin.H, 0)
	for _, t := range ticket {
		ticketResponse = append(ticketResponse, gin.H{
			"id":         t.Id,
			"ticketUuid": t.TikcetUuid,
			"ticketType": t.TicketType,
			"price":      t.Price,
			"quota":      t.Quota,
			"status":     t.Status,
			"createdAt":  t.CreatedAt,
			"updatedAt":  nil,
			"eventId":    t.EventID,
		})
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Success create ticket", ticketResponse, true))
}

// @Summary Delete tickets by ID
// @Description Deletes tickets based on their IDs
// @Tags tickets
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param ticket body dto.PayloadTicket true "List of ticket IDs to delete"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/ticket [delete]
// @Security BearerAuth
func (tc *TicketController) deleteTicket(ctx *gin.Context) {
	var payloadId dto.PayloadTicket

	if err := ctx.ShouldBindBodyWithJSON(&payloadId); err != nil {
		fmt.Println("Err JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, utils.APIResponse(err.Error(), nil, false))
		return
	}

	_, err := tc.ticketUseCase.DeleteTicketById(payloadId.Ids)
	if err != nil {
		// Perubahan disini - gunakan utils.APIResponse secara konsisten
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Success delete ticket", nil, true))
}
