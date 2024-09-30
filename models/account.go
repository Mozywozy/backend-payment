package models

import "time"

type Account struct {
    ID          uint      `gorm:"primary_key" json:"id"`
    UserID      uint      `json:"user_id"` // Foreign key ke User
    BankName    string    `json:"bank_name" binding:"required"`
    AccountNumber string  `json:"account_number" binding:"required"`
    AccountName string    `json:"account_name" binding:"required"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
