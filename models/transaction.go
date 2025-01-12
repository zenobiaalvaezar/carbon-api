package models

import (
	"time"
)

type Transaction struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	UserID        int       `json:"user_id"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	PaymentAt     time.Time `json:"payment_at"`
}
