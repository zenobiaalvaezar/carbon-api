package models

type Tree struct {
	ID             int     `gorm:"primaryKey" json:"id"`
	TreeCategoryID int     `json:"tree_category_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

//// Override default table name
//func (Category) TableName() string {
//	return "tree_category"
//}
