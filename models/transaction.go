package models

import "time"

type Transaction struct {
    ID            uint      `gorm:"primary_key" json:"id"`
    SenderID      uint      `json:"sender_id"` // ID pengguna yang mengirim
    RecipientID   uint      `json:"recipient_id"` // ID pengguna yang menerima
    Amount        float64   `json:"amount" binding:"required"`
    Type          string    `json:"type"` // Tipe transaksi, contoh: "transfer_in", "transfer_out"
    Description   string    `json:"description"`
    Status        string    `json:"status"` // e.g., "pending", "completed"
    TransactionAt time.Time `json:"transaction_at"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
