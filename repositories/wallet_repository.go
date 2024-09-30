package repositories

import (
	"payment-app/config"
	"payment-app/models"
)

type WalletRepository interface {
    GetWalletByUserID(userID uint) (*models.Wallet, error)
}

type walletRepository struct{}

func NewWalletRepository() WalletRepository {
    return &walletRepository{}
}

func (r *walletRepository) GetWalletByUserID(userID uint) (*models.Wallet, error) {
    var wallet models.Wallet
    if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
        return nil, err
    }
    return &wallet, nil
}
