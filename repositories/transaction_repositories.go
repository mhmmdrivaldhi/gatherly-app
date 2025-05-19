package repositories

import (
	"fmt"
	"gatherly-app/models"
	"gatherly-app/models/dto"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction models.Transactions) (uint, error)
	GetAll(userId int) ([]models.Transactions, error)
	FindByIdNoUser(id uint) (models.Transactions, error)
	FindById(id uint, userId int) (models.Transactions, error)
	FindByEventId(eventId uint, userId int) ([]models.Transactions, error)
	FindByTransactionIdNoUser(id string) (models.Transactions, error)
	FindByTransactionId(id string, userId int) (models.Transactions, error)
	FindByStatus(status string, userId int) ([]models.Transactions, error)
	FindByAmountRange(input dto.GetTransactionsByAmount, userId int) ([]models.Transactions, error)
	FindByDateRange(input dto.GetTransactionsByDate, userId int) ([]models.Transactions, error)
	FindByTicket(ticket string, userId int) ([]models.Transactions, error)
	DeleteById(id uint, userId int) error
	UpdateStatus(input dto.MidtransNotification)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (t *transactionRepository) Create(transaction models.Transactions) (uint, error) {
	err := t.db.Create(&transaction).Error

	if err != nil {
		return 0, err
	}

	return transaction.ID, nil
}

func (t *transactionRepository) GetAll(userId int) ([]models.Transactions, error) {
	var transactions []models.Transactions

	err := t.db.Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *transactionRepository) FindByIdNoUser(id uint) (models.Transactions, error) {
	var transaction models.Transactions

	err := t.db.First(&transaction, id).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transactionRepository) FindById(id uint, userId int) (models.Transactions, error) {
	var transaction models.Transactions

	err := t.db.Where("id = ? AND user_id = ?", id, userId).First(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transactionRepository) FindByEventId(eventId uint, userId int) ([]models.Transactions, error) {
	var transactions []models.Transactions

	err := t.db.Where("event_id = ? AND user_id = ?", eventId, userId).Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (t *transactionRepository) FindByTransactionIdNoUser(id string) (models.Transactions, error) {
	var transaction models.Transactions

	err := t.db.Where("payment_gateway_transaction_id = ?", id).First(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transactionRepository) FindByTransactionId(id string, userId int) (models.Transactions, error) {
	var transaction models.Transactions

	err := t.db.Where("payment_gateway_transaction_id = ? AND user_id = ?", id, userId).First(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transactionRepository) FindByStatus(status string, userId int) ([]models.Transactions, error) {
	var transactions []models.Transactions

	err := t.db.Where("status = ? AND user_id = ?", status, userId).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *transactionRepository) FindByAmountRange(input dto.GetTransactionsByAmount, userId int) ([]models.Transactions, error) {
	var transactions []models.Transactions

	err := t.db.Where("amount BETWEEN ? AND ? AND user_id = ?", input.MinAmount, input.MaxAmount, userId).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *transactionRepository) FindByDateRange(input dto.GetTransactionsByDate, userId int) ([]models.Transactions, error) {
	var transactions []models.Transactions

	// Define the expected date format
	layout := "2006/01/02-15:04:05"

	// Parse start and end dates
	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Execute query with parsed dates and user ID filtering
	err = t.db.Where("transaction_date BETWEEN ? AND ? AND user_id = ?", startDate, endDate, userId).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *transactionRepository) FindByTicket(ticket string, userId int) ([]models.Transactions, error) {
	var transactions []models.Transactions

	err := t.db.Where("items LIKE ? AND user_id = ?", "%"+ticket+"%", userId).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *transactionRepository) DeleteById(id uint, userId int) error {
	err := t.db.Where("id = ? AND user_id = ?", id, userId).Delete(&models.Transactions{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (t *transactionRepository) UpdateStatus(input dto.MidtransNotification) {
	t.db.Model(&models.Transactions{}).Where("payment_gateway_transaction_id = ?", input.OrderID).Updates(map[string]any{
		"status":         input.TransactionStatus,
		"payment_method": input.PaymentType,
	})

	t.db.Model(&models.EventAttendee{}).
		Where("(event_id, user_id, rsvp_date) IN (SELECT event_id, user_id, transaction_date FROM transactions WHERE payment_gateway_transaction_id = ?)", input.OrderID).
		Update("payment_status", input.TransactionStatus)
}
