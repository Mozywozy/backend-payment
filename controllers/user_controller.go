package controllers

import (
	"net/http"
	"payment-app/middlewares"
	"payment-app/models"
	"time"

	"payment-app/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashedPassword)
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func LoginUser(c *gin.Context) {
    var loginData struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := config.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Generate JWT token
    token, err := middleware.GenerateJWT(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   token,
    })
}
