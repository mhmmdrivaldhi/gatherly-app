package controllers

import (
	"errors"
	"gatherly-app/models/dto"
	"gatherly-app/usecase"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventsController struct {
	usecase usecase.EventsUsecase
	rg      *gin.RouterGroup
}

// func NewEventsController(usecase usecase.EventsUsecase, rg *gin.RouterGroup) *EventsController {
// 	return &EventsController{usecase: usecase, rg: rg}
// }

func NewEventsController(uc usecase.EventsUsecase, rg *gin.RouterGroup) *EventsController {
	return &EventsController{
		usecase: uc, // Gunakan nama parameter yang berbeda
		rg:      rg,
	}
}

func (e *EventsController) Route() {
	e.rg.POST("/event", e.createEvent)
	e.rg.GET("/event", e.getAllEvent)
	e.rg.GET("/event/:id", e.getEventByID)
	e.rg.PUT("/event/:id", e.updateEvent)
	e.rg.DELETE("/event/:id", e.deleteEvent)
	e.rg.GET("/event/distance", e.getEventByDistance)
}

// @Summary Create an event
// @Description Creates a new event
// @Tags events
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param event body dto.CreateEventRequestDTO true "Event Data"
// @Success 201 {object} dto.GeneralResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/event [post]
// @Security BearerAuth
func (e *EventsController) createEvent(ctx *gin.Context) {
	var request dto.CreateEventRequestDTO

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	createEvent, err := e.usecase.CreateEvent(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.GeneralResponse{
		Message: "successfully created event",
		Data:    createEvent,
	})
}

// @Summary Get all events
// @Description Retrieves a list of all events
// @Tags events
// @Produce json
// @Param authorization header string true "Bearer token"
// @Success 200 {object} dto.GeneralResponse
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/event [get]
// @Security BearerAuth
func (e *EventsController) getAllEvent(ctx *gin.Context) {
	events, err := e.usecase.GetAllEvent()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get all events",
		Data:    events,
	})
}

// @Summary Get event by ID
// @Description Retrieves a specific event by its ID
// @Tags events
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param id path int true "Event ID"
// @Success 200 {object} dto.GeneralResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid event ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 404 {object} dto.ErrorResponse "Event not found"
// @Router /api/v1/event/{id} [get]
// @Security BearerAuth
func (e *EventsController) getEventByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	event, err := e.usecase.GetEventByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get event by ID",
		Data:    event,
	})
}

// @Summary Update event by ID
// @Description Modifies the details of an event
// @Tags events
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param id path int true "Event ID"
// @Param event body dto.UpdateEventRequestDTO true "Updated Event Data"
// @Success 200 {object} dto.GeneralResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/event/{id} [put]
// @Security BearerAuth
func (e *EventsController) updateEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.UpdateEventRequestDTO
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	event, err := e.usecase.UpdateEvent(id, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully updated data",
		Data:    event,
	})
}

// @Summary Delete event by ID
// @Description Removes an event from the system
// @Tags events
// @Param authorization header string true "Bearer token"
// @Param id path int true "Event ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Invalid event ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 404 {object} dto.ErrorResponse "Event not found"
// @Router /api/v1/event/{id} [delete]
// @Security BearerAuth
func (e *EventsController) deleteEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = e.usecase.DeleteEvent(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, dto.GeneralResponse{
		Message: "successfully deleted data",
	})
}

// @Summary Get recommended events
// @Description Retrieves a list of recommended events for the authenticated user
// @Tags events
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param limit query int false "Number of events to return (default: 10)"
// @Success 200 {object} dto.GeneralResponse "Successfully retrieved recommended events"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameter: limit must be an integer"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/events/recommendations [get]
// @Security BearerAuth
func (e *EventsController) getEventByDistance(ctx *gin.Context) {
	latitudeValue, latitudeExists := ctx.Get("userLat")
	longitudeValue, longitudeExists := ctx.Get("userLon")

	if !latitudeExists || !longitudeExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User location not available in token"})
		return
	}

	// Konversi tipe data koordinat
	userLatitude, userLongitude, err := convertCoordinates(latitudeValue, longitudeValue)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	radiusDefault := ctx.DefaultQuery("radius", "20")
	radius, err := strconv.ParseFloat(radiusDefault, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	events, err := e.usecase.GetEventByDistance(userLatitude, userLongitude, radius)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get data by nearby location",
		Data:    events,
	})
}

// Helper function untuk konversi koordinat
func convertCoordinates(lat, lon interface{}) (float64, float64, error) {
	var userLatitude, userLongitude float64

	switch v := lat.(type) {
	case float64:
		userLatitude = v
	case float32:
		userLatitude = float64(v)
	default:
		return 0, 0, errors.New("invalid latitude type")
	}

	switch v := lon.(type) {
	case float64:
		userLongitude = v
	case float32:
		userLongitude = float64(v)
	default:
		return 0, 0, errors.New("invalid longitude type")
	}

	return userLatitude, userLongitude, nil
}
