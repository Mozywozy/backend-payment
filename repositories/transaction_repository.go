package repositories

import (
	"payment-app/config"
	"payment-app/models"
)

type TransactionRepository interface {
    Create(transaction *models.Transaction) error
}

type transactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
    return &transactionRepository{}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
    return config.DB.Create(transaction).Error
}
