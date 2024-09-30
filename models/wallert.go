package models

import "time"

type Wallet struct {
    ID        uint      `gorm:"primary_key" json:"id"`
    UserID    uint      `json:"user_id"` 
    Balance   float64   `json:"balance"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
