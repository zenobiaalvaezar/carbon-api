package models

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	RoleID    int       `json:"role_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
