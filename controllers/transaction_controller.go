package controllers

import (
	"net/http"
	"payment-app/config"
	"payment-app/models"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	var transferRequest struct {
		SenderID   uint    `json:"sender_id" binding:"required"`
		ReceiverID uint    `json:"recipient_id" binding:"required"` // Sesuaikan dengan model
		Amount     float64 `json:"amount" binding:"required"`
		Status     string  `json:"status"` // Bisa digunakan untuk menyimpan status transaksi
	}

	// Binding data request
	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi apakah pengirim dan penerima adalah pengguna yang berbeda
	if transferRequest.SenderID == transferRequest.ReceiverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender and receiver cannot be the same"})
		return
	}

	var senderWallet, receiverWallet models.Wallet

	// Dapatkan data wallet pengirim
	if err := config.DB.Where("user_id = ?", transferRequest.SenderID).First(&senderWallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender wallet not found"})
		return
	}

	// Validasi apakah saldo pengirim mencukupi
	if senderWallet.Balance < transferRequest.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Dapatkan data wallet penerima
	if err := config.DB.Where("user_id = ?", transferRequest.ReceiverID).First(&receiverWallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver wallet not found"})
		return
	}

	// Mulai transaksi database
	tx := config.DB.Begin()

	// Kurangi saldo pengirim
	senderWallet.Balance -= transferRequest.Amount
	if err := tx.Save(&senderWallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender wallet balance"})
		return
	}

	// Tambah saldo penerima
	receiverWallet.Balance += transferRequest.Amount
	if err := tx.Save(&receiverWallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receiver wallet balance"})
		return
	}

	// Buat catatan transaksi
	transaction := models.Transaction{
		SenderID:      transferRequest.SenderID,
		RecipientID:   transferRequest.ReceiverID,
		Amount:        transferRequest.Amount,
		Status:        "completed", // Status dapat disesuaikan
		TransactionAt: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Commit transaksi database
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":         "Transfer successful",
		"transaction":     transaction,
	})
}
