package models

import "encoding/json"

type Fuel struct {
	ID             int     `gorm:"primaryKey" json:"id"`
	Category       string  `json:"category"`
	Name           string  `json:"name"`
	EmissionFactor float64 `json:"emission_factor"`
	Price          float64 `json:"price"`
	Unit           string  `json:"unit"`
}

type FuelRequest struct {
	Category       string  `json:"category"`
	Name           string  `json:"name"`
	EmissionFactor float64 `json:"emission_factor"`
	Price          float64 `json:"price"`
	Unit           string  `json:"unit"`
}

func (f *Fuel) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Fuel) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}
