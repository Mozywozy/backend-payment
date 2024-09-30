package controllers

import (
	"net/http"
	"payment-app/config"
	"payment-app/models"
	"time"

	"github.com/gin-gonic/gin"
)

func AddAccount(c *gin.Context) {
    var account models.Account
    if err := c.ShouldBindJSON(&account); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    account.CreatedAt = time.Now()
    account.UpdatedAt = time.Now()

    if err := config.DB.Create(&account).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, account)
}
