package models

import "time"

type User struct {
    ID        uint      `gorm:"primary_key" json:"id"`
    Name      string    `json:"name" binding:"required"`
    Email     string    `gorm:"unique" json:"email" binding:"required,email"`
    Password  string    `json:"password"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
