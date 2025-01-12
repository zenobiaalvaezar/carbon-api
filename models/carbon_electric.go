package models

type CarbonElectric struct {
	ID               int     `gorm:"primaryKey" json:"id"`
	UserID           int     `json:"user_id"`
	ElectricID       int     `json:"electric_id"`
	Price            float64 `json:"price"`
	UsageAmount      float64 `json:"usage_amount"`
	UsageType        string  `json:"usage_type"`
	TotalConsumption float64 `json:"total_consumption"`
	EmissionFactor   float64 `json:"emission_factor"`
	EmissionAmount   float64 `json:"emission_amount"`
}
