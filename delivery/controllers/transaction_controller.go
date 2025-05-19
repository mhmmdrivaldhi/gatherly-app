package controllers

import (
	"gatherly-app/models/dto"
	"gatherly-app/usecase"
	"gatherly-app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionUsecase usecase.TransactionUsecase
	rg                 *gin.RouterGroup
}

func NewTransactionController(transactionUsecase usecase.TransactionUsecase, rg *gin.RouterGroup) *TransactionController {
	return &TransactionController{transactionUsecase: transactionUsecase, rg: rg}
}

func (t *TransactionController) Route() {
	t.rg.GET("/transactions", t.getAllTransactions)
	t.rg.GET("/transaction/:id", t.findTransactionById)
	t.rg.GET("/transaction/event-id/:event_id", t.findTransactionByEventId)
	t.rg.GET("/transaction/transaction-id/:id", t.findTransactionByTransactionId)
	t.rg.GET("/transaction/status/:status", t.findTransactionByStatus)
	t.rg.POST("/transaction/amount-range", t.findTransactionByAmountRange)
	t.rg.POST("/transaction/date-range", t.findTransactionByDateRange)
	t.rg.GET("/transaction/ticket/:ticket", t.findTransactionByTicket)
}

func (t *TransactionController) RegisterPublicRoutes() {
	t.rg.POST("/transaction/notification", t.handleNotification)
}

// @Summary Get all transactions
// @Description Retrieves all transactions
// @Tags transactions
// @Produce json
// @Param authorization header string true "Bearer token"
// @Success 200 {object} utils.Response
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transactions [get]
// @Security BearerAuth
func (t *TransactionController) getAllTransactions(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	result, err := t.transactionUsecase.GetAllTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, utils.APIResponse("No transactions found", nil, false))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success get all transactions", result, true))
}

// @Summary Get transaction by ID
// @Description Retrieves a specific transaction
// @Tags transactions
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param id path int true "Transaction ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid transaction ID"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/{id} [get]
// @Security BearerAuth
func (t *TransactionController) findTransactionById(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id := uint(id64)
	result, err := t.transactionUsecase.FindTransactionById(id, userID)
	if err != nil && err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"err": "Transaction not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success find transaction", result, true))
}

// @Summary Get transactions by Event ID
// @Description Retrieves all transactions linked to a specific event
// @Tags transactions
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param event_id path int true "Event ID"
// @Success 200 {array} utils.Response
// @Failure 400 {object} string "Invalid event ID"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/event-id/{event_id} [get]
// @Security BearerAuth
func (t *TransactionController) findTransactionByEventId(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	eventID64, err := strconv.ParseUint(c.Param("event_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	eventID := uint(eventID64)
	result, err := t.transactionUsecase.FindTransactionByEventId(eventID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, utils.APIResponse("No transactions found", nil, false))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success find transactions", result, true))
}

// @Summary Find transaction by Transaction ID
// @Description Retrieves a transaction using its transaction ID
// @Tags transactions
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param id path string true "Transaction ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid transaction ID"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/transaction-id/{id} [get]
// @Security BearerAuth
func (t *TransactionController) findTransactionByTransactionId(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	id := c.Param("id")
	result, err := t.transactionUsecase.FindTransactionByTransactionId(id, userID)
	if err != nil && err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"err": "Transaction not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.APIResponse("Success find transaction", result, true))
}

// @Summary Find transactions by Status
// @Description Retrieves all transactions with a specific status
// @Tags transactions
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param status path string true "Transaction Status"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid status"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/status/{status} [get]
// @Security BearerAuth
func (t *TransactionController) findTransactionByStatus(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	status := c.Param("status")
	result, err := t.transactionUsecase.FindTransactionByStatus(status, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, utils.APIResponse("No transactions found", nil, false))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success find transactions by status", result, true))
}

// @Summary Find transactions by Amount Range
// @Description Retrieves transactions within a given amount range
// @Tags transactions
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param request body dto.GetTransactionsByAmount true "Amount Range Payload"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/amount-range [post]
// @Security BearerAuth
func (t *TransactionController) findTransactionByAmountRange(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	var input dto.GetTransactionsByAmount
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	result, err := t.transactionUsecase.FindTransactionByAmountRange(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, utils.APIResponse("No transactions found", nil, false))
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, utils.APIResponse("No transactions found", nil, false))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success find transactions by amount range", result, true))
}

// @Summary Find transactions by Date Range
// @Description Retrieves transactions within a specific date range
// @Tags transactions
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param request body dto.GetTransactionsByDate true "Date Range Payload"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/date-range [post]
// @Security BearerAuth
func (t *TransactionController) findTransactionByDateRange(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	var input dto.GetTransactionsByDate
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	result, err := t.transactionUsecase.FindTransactionByDateRange(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, utils.APIResponse("No transactions found", nil, false))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success find transactions by date range", result, true))
}

// @Summary Find transactions by Ticket
// @Description Retrieves transactions by ticket information
// @Tags transactions
// @Produce json
// @Param authorization header string true "Bearer token"
// @Param ticket path string true "Ticket"
// @Success 200 {object} utils.Response
// @Failure 400 {object} string "Invalid ticket"
// @Failure 401 {object} string "Unauthorized: Missing or invalid token"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/ticket/{ticket} [get]
// @Security BearerAuth
func (t *TransactionController) findTransactionByTicket(c *gin.Context) {
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.APIResponse("User ID not found in token", nil, false))
		return
	}

	userID, ok := userIDFromToken.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.APIResponse("User ID in token is of invalid type", nil, false))
		return
	}

	ticket := c.Param("ticket")
	result, err := t.transactionUsecase.FindTransactionByTicket(ticket, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, utils.APIResponse("No transactions found", nil, false))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success find transactions by ticket", result, true))
}

// @Summary Handle payment notification
// @Description Processes Midtrans payment notifications and updates transaction status
// @Tags transactions
// @Accept json
// @Produce json
// @Param notification body dto.MidtransNotification true "Midtrans Payment Notification"
// @Success 200 {object} utils.Response "Success handle notification"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/transaction/notification [post]
func (t *TransactionController) handleNotification(c *gin.Context) {
	var notification dto.MidtransNotification
	err := c.ShouldBindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = t.transactionUsecase.HandleNotification(notification)
	if err != nil && err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"err": "Transaction not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.APIResponse("Success handle notification", nil, true))
}
