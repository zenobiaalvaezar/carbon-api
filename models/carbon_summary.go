package models

type CarbonSummary struct {
	ID               int     `gorm:"primaryKey" json:"id"`
	UserID           int     `json:"user_id"`
	FuelEmission     float64 `json:"fuel_emission"`
	ElectricEmission float64 `json:"electric_emission"`
	TotalEmission    float64 `json:"total_emission"`
	TotalTree        int     `json:"total_tree"`
}
