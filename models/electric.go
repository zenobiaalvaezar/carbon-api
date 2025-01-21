package models

type Electric struct {
	ID             int     `gorm:"primaryKey" json:"id"`
	Province       string  `json:"province"`
	EmissionFactor float64 `json:"emission_factor"`
	Price          float64 `json:"price"`
}
