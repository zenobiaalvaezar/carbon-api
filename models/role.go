package models

type Role struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
