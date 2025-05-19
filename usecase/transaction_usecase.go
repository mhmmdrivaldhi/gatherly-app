package usecase

import (
	"fmt"
	"gatherly-app/models"
	"gatherly-app/models/dto"
	"gatherly-app/repositories"
	"gatherly-app/service"

	"time"
)

type TransactionUsecase interface {
	CreateTransaction(input dto.CreateTransaction, midtransInput dto.MidtransSnapReq) (models.Transactions, error)
	GetAllTransactions(userId int) ([]models.Transactions, error)
	FindTransactionById(id uint, userId int) (models.Transactions, error)
	FindTransactionByEventId(eventId uint, userId int) ([]models.Transactions, error)
	FindTransactionByTransactionId(id string, userId int) (models.Transactions, error)
	FindTransactionByStatus(status string, userId int) ([]models.Transactions, error)
	FindTransactionByAmountRange(input dto.GetTransactionsByAmount, userId int) ([]models.Transactions, error)
	FindTransactionByDateRange(input dto.GetTransactionsByDate, userId int) ([]models.Transactions, error)
	FindTransactionByTicket(ticket string, userId int) ([]models.Transactions, error)
	DeleteTransactionById(id uint, userId int) error
	HandleNotification(notification dto.MidtransNotification) error
}

type transactionUsecase struct {
	transactionRepository repositories.TransactionRepository
	midtransService       service.MidtransService
}

func NewTransactionUsecase(transactionRepository repositories.TransactionRepository, midtransService service.MidtransService) TransactionUsecase {
	return &transactionUsecase{
		transactionRepository: transactionRepository,
		midtransService:       midtransService,
	}
}

func (t *transactionUsecase) CreateTransaction(input dto.CreateTransaction, midtransInput dto.MidtransSnapReq) (models.Transactions, error) {
	resp, err := t.midtransService.Pay(midtransInput)
	if err != nil {
		return models.Transactions{}, err
	}

	transaction := models.Transactions{
		UserId:                      input.UserId,
		EventId:                     input.EventId,
		Amount:                      input.Amount,
		TransactionDate:             time.Now(),
		Status:                      "pending",
		PaymentMethod:               "",
		PaymentGatewayTransactionId: midtransInput.TransactionDetails.OrderID,
		Items:                       input.Items,
		Notes:                       input.Notes,
		Url:                         resp.RedirectURL,
	}

	id, err := t.transactionRepository.Create(transaction)

	if err != nil {
		return models.Transactions{}, err
	}

	result, err := t.transactionRepository.FindByIdNoUser(id)

	if err != nil {
		return models.Transactions{}, err
	}

	return result, nil
}

func (t *transactionUsecase) GetAllTransactions(userId int) ([]models.Transactions, error) {
	result, err := t.transactionRepository.GetAll(userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *transactionUsecase) FindTransactionById(id uint, userId int) (models.Transactions, error) {
	result, err := t.transactionRepository.FindById(id, userId)

	if err != nil {
		return models.Transactions{}, err
	}

	return result, nil
}

func (t *transactionUsecase) FindTransactionByEventId(eventId uint, userId int) ([]models.Transactions, error) {
	result, err := t.transactionRepository.FindByEventId(eventId, userId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *transactionUsecase) FindTransactionByTransactionId(id string, userId int) (models.Transactions, error) {
	result, err := t.transactionRepository.FindByTransactionId(id, userId)

	if err != nil {
		return models.Transactions{}, err
	}

	return result, nil
}

func (t *transactionUsecase) FindTransactionByStatus(status string, userId int) ([]models.Transactions, error) {
	result, err := t.transactionRepository.FindByStatus(status, userId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *transactionUsecase) FindTransactionByAmountRange(input dto.GetTransactionsByAmount, userId int) ([]models.Transactions, error) {
	result, err := t.transactionRepository.FindByAmountRange(input, userId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *transactionUsecase) FindTransactionByDateRange(input dto.GetTransactionsByDate, userId int) ([]models.Transactions, error) {
	result, err := t.transactionRepository.FindByDateRange(input, userId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *transactionUsecase) FindTransactionByTicket(ticket string, userId int) ([]models.Transactions, error) {
	result, err := t.transactionRepository.FindByTicket(ticket, userId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *transactionUsecase) DeleteTransactionById(id uint, userId int) error {
	// Check if the transaction exists and belongs to the user
	transaction, err := t.transactionRepository.FindById(id, userId)
	if err != nil {
		return err
	}

	// Delete only if status is "pending"
	if transaction.Status != "pending" {
		return fmt.Errorf("cannot delete transaction with status: %s", transaction.Status)
	}

	err = t.transactionRepository.DeleteById(id, userId)
	if err != nil {
		return err
	}

	return nil
}

func (t *transactionUsecase) HandleNotification(notification dto.MidtransNotification) error {
	_, err := t.transactionRepository.FindByTransactionIdNoUser(notification.OrderID)

	if err != nil {
		return err
	}

	t.transactionRepository.UpdateStatus(notification)
	return nil
}
