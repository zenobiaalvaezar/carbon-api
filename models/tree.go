package models

type Tree struct {
	ID             int     `gorm:"primaryKey" json:"id"`
	TreeCategoryID int     `json:"tree_category_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
}
