package models

type TransactionDetail struct {
	ID            int     `gorm:"primaryKey" json:"id"`
	TransactionID int     `json:"transaction_id"`
	TreeID        int     `json:"tree_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	TotalPrice    float64 `json:"total_price"`
}
