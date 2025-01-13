package models

type Cart struct {
	ID       int `gorm:"primaryKey" json:"id"`
	UserID   int `json:"user_id"`
	TreeID   int `json:"tree_id"`
	Quantity int `json:"quantity"`
}
