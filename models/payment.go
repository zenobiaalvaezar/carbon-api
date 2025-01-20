package models

type PaymentRequest struct {
	TransactionID int     `json:"transaction_id"`
	UserID        int     `json:"user_id"`
	PaymentMethod string  `json:"payment_method"`
	PaymentAmount float64 `json:"payment_amount"`
}

type PaymentResponse struct {
	TransactionID int     `json:"transaction_id"`
	PaymentMethod string  `json:"payment_method"`
	PaymentAmount float64 `json:"payment_amount"`
	PaymentStatus string  `json:"payment_status"`
	RedirectURL   string  `json:"redirect_url"`
}

type InvoiceRequest struct {
	ExternalId         string  `json:"external_id"`
	Amount             float64 `json:"amount"`
	Description        string  `json:"description"`
	InvoiceDuration    int     `json:"invoice_duration"`
	GivenNames         string  `json:"given_names"`
	Email              string  `json:"email"`
	Currency           string  `json:"currency"`
	PaymentMethod      string  `json:"payment_method"`
	SuccessRedirectURL string  `json:"success_redirect_url"`
	FailureRedirectURL string  `json:"failure_redirect_url"`
}

type InvoiceResponse struct {
	Id         string  `json:"id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
	InvoiceURL string  `json:"invoice_url"`
}
