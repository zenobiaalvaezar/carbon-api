package models

import "time"

type ReportSummary struct {
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	CarbonTree   int    `json:"carbon_tree"`
	DonationTree int    `json:"donation_tree"`
	BadgeStatus  string `json:"badge_status"`
}

type ReportDetail struct {
	UserID          int       `json:"user_id"`
	UserName        string    `json:"user_name"`
	UserEmail       string    `json:"user_email"`
	TransactionID   int       `json:"transaction_id"`
	TransactionDate time.Time `json:"transaction_date"`
	TotalTree       int       `json:"total_tree"`
	TotalPrice      float64   `json:"total_price"`
	PaymentMethod   string    `json:"payment_method"`
	PaymentStatus   string    `json:"payment_status"`
	PaymentAt       time.Time `json:"payment_at"`
}
