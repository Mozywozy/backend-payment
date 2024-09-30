package services

import (
	"payment-app/models"
	"payment-app/repositories"
)

type WalletService interface {
    GetWallet(userID uint) (*models.Wallet, error)
}

type walletService struct {
    walletRepository repositories.WalletRepository
}

func NewWalletService(walletRepo repositories.WalletRepository) WalletService {
    return &walletService{
        walletRepository: walletRepo,
    }
}

func (s *walletService) GetWallet(userID uint) (*models.Wallet, error) {
    return s.walletRepository.GetWalletByUserID(userID)
}
