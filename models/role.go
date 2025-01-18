package models

type Role struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type RoleRequest struct {
	Name string `json:"name"`
}
