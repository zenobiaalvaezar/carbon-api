package models

type ReportSummary struct {
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	CarbonTree   int    `json:"carbon_tree"`
	DonationTree int    `json:"donation_tree"`
	BadgeStatus  string `json:"badge_status"`
}
