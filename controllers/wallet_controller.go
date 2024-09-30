package controllers

import (
	"net/http"
	"payment-app/config"
	"payment-app/models"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateWallet(c *gin.Context) {
	var walletRequest struct {
		UserID uint    `json:"user_id" binding:"required"`
		Balance float64 `json:"balance" binding:"required"`
	}

	// Bind data dari request
	if err := c.ShouldBindJSON(&walletRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Membuat objek wallet baru
	newWallet := models.Wallet{
		UserID:    walletRequest.UserID,
		Balance:   walletRequest.Balance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Simpan wallet ke database
	if err := config.DB.Create(&newWallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Wallet created successfully",
		"wallet":  newWallet,
	})
}

func GetWallet(c *gin.Context) {
    userID := c.Param("userID")

    var wallet models.Wallet
    if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
        return
    }

    c.JSON(http.StatusOK, wallet)
}

// AddBalance untuk menambah saldo dompet pengguna
func AddBalance(c *gin.Context) {
	var addBalanceRequest struct {
		UserID uint    `json:"user_id" binding:"required"`
		Amount  float64 `json:"amount" binding:"required"`
	}

	// Bind data dari request
	if err := c.ShouldBindJSON(&addBalanceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wallet models.Wallet

	// Dapatkan data wallet pengguna
	if err := config.DB.Where("user_id = ?", addBalanceRequest.UserID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Tambah saldo
	wallet.Balance += addBalanceRequest.Amount
	wallet.UpdatedAt = time.Now()

	if err := config.DB.Save(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Balance updated successfully",
		"wallet":  wallet,
	})
}