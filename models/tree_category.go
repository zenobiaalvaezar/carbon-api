package models

type TreeCategory struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
