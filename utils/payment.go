package utils

import (
	"bytes"
	"carbon-api/models"
	"encoding/json"
	"net/http"
	"os"
)

func CreateInvoice(req models.InvoiceRequest) (*models.InvoiceResponse, int, error) {
	apiURL := os.Getenv("XENDIT_API_URL")
	apiKey := os.Getenv("XENDIT_SECRET_KEY")

	body := map[string]interface{}{
		"external_id":      req.ExternalId,
		"amount":           req.Amount,
		"description":      req.Description,
		"invoice_duration": req.InvoiceDuration,
		"customer": map[string]string{
			"given_names": req.GivenNames,
			"email":       req.Email,
		},
		"currency":             req.Currency,
		"payment_methods":      []string{req.PaymentMethod},
		"success_redirect_url": req.SuccessRedirectURL,
		"failure_redirect_url": req.FailureRedirectURL,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer response.Body.Close()

	var resInvoice models.InvoiceResponse
	if err := json.NewDecoder(response.Body).Decode(&resInvoice); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &resInvoice, response.StatusCode, nil
}
