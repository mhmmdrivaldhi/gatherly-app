package controllers

import (
	"gatherly-app/models/dto"
	"gatherly-app/usecase"
	"gatherly-app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventAttendeeController struct {
	eventAttendeeUseCase usecase.EventAttendeeUseCase
	rg                   *gin.RouterGroup
}

func NewEventAttendeeController(eventAttendeeUseCase usecase.EventAttendeeUseCase, rg *gin.RouterGroup) *EventAttendeeController {
	return &EventAttendeeController{
		eventAttendeeUseCase: eventAttendeeUseCase,
		rg:                   rg,
	}
}

func (ec *EventAttendeeController) Route() {
	ec.rg.POST("/attendee", ec.Register)
	ec.rg.DELETE("/attendee", ec.Cancel)
	ec.rg.GET("/attendee", ec.GetRegistration)
	ec.rg.GET("/attendee/event/:eventId", ec.ListAttendeesByEvent)
	ec.rg.GET("/attendee/user/:userId", ec.ListRegistrationsByUser)
	ec.rg.PATCH("/attendee/confirm-payment", ec.ConfirmPayment)
	ec.rg.PATCH("/attendee/rsvp", ec.UpdateRSVP)
}

// @Summary Register for an event
// @Description Adds an attendee to an event
// @Tags event_attendees
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param registration body dto.AttendeeRegisterRequest true "Attendee Registration Data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/attendee [post]
// @Security BearerAuth
func (ec *EventAttendeeController) Register(ctx *gin.Context) {
	var payload dto.AttendeeRegisterRequest // AttendeeRegisterRequest is in fp_gatherly/model/dto/event_attendees_dto.go
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.APIResponse(err.Error(), nil, false))
		return
	}

	// Get userID from context set by auth middleware
	userIDFromToken, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	// Use userID from token instead of payload.UserID
	attendee, err := ec.eventAttendeeUseCase.Register(ctx, userID, payload.EventID, payload.TicketTypeID, payload.RSVPStatus) // Added payload.TicketTypeID
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Success register for event", attendee, true))
}

// @Summary Cancel event registration
// @Description Removes an attendee from an event
// @Tags event_attendees
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param cancellation body dto.AttendeeCancelRequest true "Attendee Cancellation Data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/attendee [delete]
// @Security BearerAuth
func (ec *EventAttendeeController) Cancel(ctx *gin.Context) {
	var payload dto.AttendeeCancelRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.APIResponse(err.Error(), nil, false))
		return
	}

	err := ec.eventAttendeeUseCase.CancelRegistration(ctx, payload.UserID, payload.EventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Successfully cancelled registration", nil, true))
}

// @Summary Get attendee registration details
// @Description Retrieves registration details for a specific user and event
// @Tags event_attendees
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param userId query int true "User ID"
// @Param eventId query int true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid query parameters"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/attendee [get]
// @Security BearerAuth
func (ec *EventAttendeeController) GetRegistration(ctx *gin.Context) {
	userIDStr := ctx.Query("userId")
	eventIDStr := ctx.Query("eventId")

	userID, err1 := strconv.Atoi(userIDStr)
	eventID, err2 := strconv.Atoi(eventIDStr)

	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, utils.APIResponse("Invalid query parameters", nil, false))
		return
	}

	attendee, err := ec.eventAttendeeUseCase.GetRegistrationDetails(ctx, userID, eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Success fetch registration details", attendee, true))
}

// @Summary List attendees of an event
// @Description Retrieves all attendees for a specific event
// @Tags event_attendees
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param eventId path int true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid event ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/attendee/event/{eventId} [get]
// @Security BearerAuth
func (ec *EventAttendeeController) ListAttendeesByEvent(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Param("eventId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.APIResponse("Invalid event ID", nil, false))
		return
	}

	attendees, err := ec.eventAttendeeUseCase.ListAttendeesForEvent(ctx, eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Success fetch attendees", attendees, true))
}

// @Summary List a user's event registrations
// @Description Retrieves all event registrations for a specific user
// @Tags event_attendees
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param userId path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid user ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/attendee/user/{userId} [get]
// @Security BearerAuth
func (ec *EventAttendeeController) ListRegistrationsByUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.APIResponse("Invalid user ID", nil, false))
		return
	}

	registrations, err := ec.eventAttendeeUseCase.ListUserRegistrations(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Success fetch user registrations", registrations, true))
}

// @Summary Confirm payment for an event
// @Description Updates payment status for an event registration
// @Tags event_attendees
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param payment body dto.AttendeePaymentRequest true "Payment confirmation data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/attendee/confirm-payment [patch]
// @Security BearerAuth
func (ec *EventAttendeeController) ConfirmPayment(ctx *gin.Context) {
	var payload dto.AttendeePaymentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.APIResponse(err.Error(), nil, false))
		return
	}

	attendee, err := ec.eventAttendeeUseCase.ConfirmPayment(ctx, payload.UserID, payload.EventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("Payment confirmed", attendee, true))
}

// @Summary Update RSVP status
// @Description Updates an attendee's RSVP status for an event
// @Tags event_attendees
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param rsvp body dto.AttendeeRSVPRequest true "RSVP Update Data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/attendee/rsvp [patch]
// @Security BearerAuth
func (ec *EventAttendeeController) UpdateRSVP(ctx *gin.Context) {
	var payload dto.AttendeeRSVPRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.APIResponse(err.Error(), nil, false))
		return
	}

	attendee, err := ec.eventAttendeeUseCase.UpdateRSVPStatus(ctx, payload.UserID, payload.EventID, payload.NewStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, utils.APIResponse("RSVP status updated", attendee, true))
}
