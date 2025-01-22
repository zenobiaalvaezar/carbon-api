package repositories

import (
	"carbon-api/models"
	"carbon-api/utils"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment models.PaymentRequest) (models.PaymentResponse, int, error)
	VerifyPayment(transactionID int, status string) (int, error)
}

type paymentRepository struct {
	DB              *gorm.DB
	MongoCollection *mongo.Collection
}

func NewPaymentRepository(DB *gorm.DB, MongoCollection *mongo.Collection) *paymentRepository {
	return &paymentRepository{DB, MongoCollection}
}

func (r *paymentRepository) CreatePayment(payment models.PaymentRequest) (models.PaymentResponse, int, error) {
	var transaction models.Transaction
	r.DB.Where("id = ?", payment.TransactionID).First(&transaction)
	if transaction.ID == 0 {
		return models.PaymentResponse{}, http.StatusNotFound, errors.New("Transaction not found")
	}

	var user models.User
	r.DB.Where("id = ?", transaction.UserID).First(&user)
	if user.ID == 0 {
		return models.PaymentResponse{}, http.StatusNotFound, errors.New("User not found")
	}

	// check payment method is exists in payment method list
	paymentMethodRepository := NewPaymentMethodRepository(r.MongoCollection)
	_, _, err := paymentMethodRepository.GetPaymentMethodByCode(payment.PaymentMethod)
	if err != nil {
		return models.PaymentResponse{}, http.StatusNotFound, err
	}

	transaction.PaymentMethod = payment.PaymentMethod
	transaction.PaymentStatus = "pending"
	r.DB.Save(&transaction)

	// create invoice
	baseURL := os.Getenv("BASE_URL")
	invoice := models.InvoiceRequest{
		ExternalId:         strconv.Itoa(payment.TransactionID),
		Amount:             transaction.TotalPrice,
		Description:        "Top up deposit",
		InvoiceDuration:    86400,
		GivenNames:         user.Name,
		Email:              user.Email,
		Currency:           "IDR",
		PaymentMethod:      payment.PaymentMethod,
		SuccessRedirectURL: baseURL + "/payments/verify/" + strconv.Itoa(payment.TransactionID) + "?status=success",
		FailureRedirectURL: baseURL + "/payments/verify/" + strconv.Itoa(payment.TransactionID) + "?status=failed",
	}

	// create invoice
	resInvoice, statusCode, err := utils.CreateInvoice(invoice)
	if err != nil {
		return models.PaymentResponse{}, statusCode, err
	}

	return models.PaymentResponse{
		TransactionID: transaction.ID,
		PaymentMethod: transaction.PaymentMethod,
		PaymentAmount: transaction.TotalPrice,
		PaymentStatus: transaction.PaymentStatus,
		RedirectURL:   resInvoice.InvoiceURL,
	}, http.StatusOK, nil
}

func (r *paymentRepository) VerifyPayment(transactionID int, status string) (int, error) {
	var transaction models.Transaction
	r.DB.Where("id = ?", transactionID).First(&transaction)
	if transaction.ID == 0 {
		return http.StatusNotFound, errors.New("Transaction not found")
	}

	// check if payment status is not pending
	if transaction.PaymentStatus != "pending" {
		return http.StatusBadRequest, errors.New("Transaction is not pending")
	}

	// if payment failed, return stock to tree
	if status == "failed" {
		var details []models.TransactionDetail
		r.DB.Where("transaction_id = ?", transactionID).Find(&details)

		for _, detail := range details {
			var tree models.Tree
			r.DB.Where("id = ?", detail.TreeID).First(&tree)
			tree.Stock += detail.Quantity
			r.DB.Save(&tree)
		}

		transaction.PaymentStatus = "failed"
		r.DB.Save(&transaction)
	} else {
		transaction.PaymentStatus = "success"
		transaction.PaymentAt = time.Now()
		r.DB.Save(&transaction)
	}

	return http.StatusOK, nil
}
