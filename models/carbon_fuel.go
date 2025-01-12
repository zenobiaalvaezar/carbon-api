package models

type CarbonFuel struct {
	ID               int     `gorm:"primaryKey" json:"id"`
	UserID           int     `json:"user_id"`
	FuelID           int     `json:"fuel_id"`
	Price            float64 `json:"price"`
	UsageAmount      float64 `json:"usage_amount"`
	UsageType        string  `json:"usage_type"`
	TotalConsumption float64 `json:"total_consumption"`
	EmissionFactor   float64 `json:"emission_factor"`
	EmissionAmount   float64 `json:"emission_amount"`
}

type CarbonFuelRequest struct {
	UserID      int     `json:"user_id"`
	FuelID      int     `json:"fuel_id"`
	UsageType   string  `json:"usage_type"`
	UsageAmount float64 `json:"usage_amount"`
}

type CarbonFuelResponse struct {
	ID               int     `json:"id"`
	UserID           int     `json:"user_id"`
	UserName         string  `json:"user_name"`
	UserEmail        string  `json:"user_email"`
	FuelID           int     `json:"fuel_id"`
	FuelName         string  `json:"fuel_name"`
	Price            float64 `json:"price"`
	Unit             string  `json:"unit"`
	UsageAmount      float64 `json:"usage_amount"`
	UsageType        string  `json:"usage_type"`
	TotalConsumption float64 `json:"total_consumption"`
	EmissionFactor   float64 `json:"emission_factor"`
	EmissionAmount   float64 `json:"emission_amount"`
}
