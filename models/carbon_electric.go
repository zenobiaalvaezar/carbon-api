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

type CarbonElectricResponse struct {
	ID               int     `json:"id"`
	UserID           int     `json:"user_id"`
	UserName         string  `json:"user_name"`
	UserEmail        string  `json:"user_email"`
	UsageType        string  `json:"usage_type"`
	UsageAmount      float64 `json:"usage_amount"`
	TotalConsumption float64 `json:"total_consumption"`
	EmissionFactor   float64 `json:"emission_factor"`
	EmissionAmount   float64 `json:"emission_amount"`
}

type CarbonElectricRequest struct {
	UserID         int     `json:"user_id" validate:"required"`           
	UsageType      string  `json:"usage_type" validate:"required,oneof=consumption rupiah"` 
	UsageAmount    float64 `json:"usage_amount" validate:"required,gt=0"` 
	Price          float64 `json:"price" validate:"required,gt=0"`         
	EmissionFactor float64 `json:"emission_factor" validate:"required,gt=0"` 
}