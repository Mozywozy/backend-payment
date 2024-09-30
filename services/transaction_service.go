package services

import (
	"payment-app/models"
	"payment-app/repositories"
)

type TransactionService interface {
    CreateTransaction(transaction *models.Transaction) error
}

type transactionService struct {
    transactionRepository repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) TransactionService {
    return &transactionService{
        transactionRepository: transactionRepo,
    }
}

func (s *transactionService) CreateTransaction(transaction *models.Transaction) error {
    return s.transactionRepository.Create(transaction)
}
