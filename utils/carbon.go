package utils

import "math"

func CalculateTotalConsumption(usageType string, usageAmount float64, price float64) float64 {
	var totalConsumption float64
	if usageType == "rupiah" {
		totalConsumption = usageAmount / price
	} else {
		totalConsumption = usageAmount
	}
	return math.Round(totalConsumption*100) / 100
}

func CalculateEmissionAmount(totalConsumption float64, emissionFactor float64) float64 {
	return math.Round((totalConsumption*emissionFactor)*100) / 100
}

func CalculateTotalEmission(fuelEmission float64, electricEmission float64) float64 {
	return math.Round((fuelEmission+electricEmission)*100) / 100
}

func CalculateTotalTree(totalEmission float64) int {
	return int(math.Round(totalEmission / 34.5))
}
