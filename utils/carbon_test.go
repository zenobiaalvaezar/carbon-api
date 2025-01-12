package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTotalConsumption(t *testing.T) {
	type args struct {
		usageType   string
		usageAmount float64
		price       float64
	}
	type test struct {
		name     string
		args     args
		expected float64
	}

	tests := []test{
		{
			name: "Test Calculate Total Consumption",
			args: args{
				usageType:   "rupiah",
				usageAmount: 100000,
				price:       10000,
			},
			expected: 10,
		},
		{
			name: "Test Calculate Total Consumption",
			args: args{
				usageType:   "consumption",
				usageAmount: 10,
				price:       10000,
			},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalculateTotalConsumption(tt.args.usageType, tt.args.usageAmount, tt.args.price)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCalculateEmissionAmount(t *testing.T) {
	type args struct {
		totalConsumption float64
		emissionFactor   float64
	}
	type test struct {
		name     string
		args     args
		expected float64
	}

	tests := []test{
		{
			name: "Test Calculate Emission Amount",
			args: args{
				totalConsumption: 10,
				emissionFactor:   10,
			},
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalculateEmissionAmount(tt.args.totalConsumption, tt.args.emissionFactor)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCalculateTotalEmission(t *testing.T) {
	type args struct {
		fuelEmission     float64
		electricEmission float64
	}
	type test struct {
		name     string
		args     args
		expected float64
	}

	tests := []test{
		{
			name: "Test Calculate Total Emission",
			args: args{
				fuelEmission:     10,
				electricEmission: 10,
			},
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalculateTotalEmission(tt.args.fuelEmission, tt.args.electricEmission)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCalculateTotalTree(t *testing.T) {
	// 1 tree can absorb 34.5 pounds of CO2 per year
	type args struct {
		totalEmission float64
	}
	type test struct {
		name     string
		args     args
		expected int
	}

	tests := []test{
		{
			name: "Test Calculate Total Tree",
			args: args{
				totalEmission: 34.5,
			},
			expected: 1,
		},
		{
			name: "Test Calculate Total Tree",
			args: args{
				totalEmission: 69,
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalculateTotalTree(tt.args.totalEmission)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
