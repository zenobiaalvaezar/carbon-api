package models

type Cart struct {
	ID       int `gorm:"primaryKey" json:"id"`
	UserID   int `json:"user_id"`
	TreeID   int `json:"tree_id"`
	Quantity int `json:"quantity"`
}

type AddCartRequest struct {
	UserID   int `json:"user_id"`
	TreeID   int `json:"tree_id"`
	Quantity int `json:"quantity"`
}

type GetCartsResponse struct {
	ID       int     `json:"id"`
	TreeName string  `json:"tree_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
