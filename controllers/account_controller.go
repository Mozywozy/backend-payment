package controllers

import (
	"net/http"
	"payment-app/config"
	"payment-app/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AddAccount(c *gin.Context) {
    var accountRequest models.Account

    // Bind data request
    if err := c.ShouldBindJSON(&accountRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash the PIN
    hashedPIN, err := bcrypt.GenerateFromPassword([]byte(accountRequest.PIN), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash PIN"})
        return
    }

    // Set hashed PIN
    accountRequest.PIN = string(hashedPIN)

    // Save to the database
    if err := config.DB.Create(&accountRequest).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Account created successfully", "account": accountRequest})
}