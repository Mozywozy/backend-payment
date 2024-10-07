package controllers

import (
	"net/http"
	"payment-app/config"
	"payment-app/models"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
    var transferRequest struct {
        SenderID   uint    `json:"sender_id" binding:"required"`
        ReceiverID uint    `json:"receiver_id" binding:"required"`
        Amount     float64 `json:"amount" binding:"required"`
        PIN        string  `json:"pin" binding:"required"`
        Description string `json:"description"`
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

    var senderWallet models.Wallet
    var senderAccount models.Account
    var receiverWallet models.Wallet

    // Dapatkan data akun dan wallet pengirim
    if err := config.DB.Where("user_id = ?", transferRequest.SenderID).First(&senderAccount).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Sender account not found"})
        return
    }

    // Validasi PIN pengirim
	if err := bcrypt.CompareHashAndPassword([]byte(senderAccount.PIN), []byte(transferRequest.PIN)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PIN"})
		return
	}

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

    // Buat catatan transaksi untuk pengirim
    senderTransaction := models.Transaction{
        SenderID:        transferRequest.SenderID,
        Amount:        -transferRequest.Amount,
        Type:          "transfer_out",
        Description:   transferRequest.Description,
        Status:        "completed",
        TransactionAt: time.Now(),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if err := tx.Create(&senderTransaction).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sender transaction"})
        return
    }

    // Buat catatan transaksi untuk penerima
    receiverTransaction := models.Transaction{
        SenderID:        transferRequest.ReceiverID,
        Amount:        transferRequest.Amount,
        Type:          "transfer_in",
        Description:   transferRequest.Description,
        Status:        "completed",
        TransactionAt: time.Now(),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if err := tx.Create(&receiverTransaction).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create receiver transaction"})
        return
    }

    // Commit transaksi database
    tx.Commit()

    c.JSON(http.StatusOK, gin.H{
        "message":            "Transfer successful",
        "sender_transaction": senderTransaction,
        "receiver_transaction": receiverTransaction,
    })
}


func GetTransactionHistory(c *gin.Context) {
    userID := c.Param("userID") 

    var transactions []models.Transaction

    if err := config.DB.Where("sender_id = ?", userID).Order("created_at desc").Find(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction history"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}