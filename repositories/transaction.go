package repositories

import (
	"carbon-api/models"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetAllTransactions(userId int) ([]models.TransactionResponse, int, error)
	AddTransaction(userId int) (models.Transaction, int, error)
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) *transactionRepository {
	return &transactionRepository{DB}
}

func (r *transactionRepository) GetAllTransactions(userId int) ([]models.TransactionResponse, int, error) {
	var transactions []models.Transaction
	err := r.DB.Where("user_id = ?", userId).Find(&transactions).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var response []models.TransactionResponse
	for _, transaction := range transactions {
		var details []models.TransactionDetail
		err := r.DB.Where("transaction_id = ?", transaction.ID).Find(&details).Error
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var detailResponse []models.TransactionDetailResponse
		for _, detail := range details {
			var tree models.Tree
			err := r.DB.Where("id = ?", detail.TreeID).First(&tree).Error
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}

			detailResponse = append(detailResponse, models.TransactionDetailResponse{
				ID:         detail.ID,
				TreeName:   tree.Name,
				Quantity:   detail.Quantity,
				Price:      detail.Price,
				TotalPrice: detail.TotalPrice,
			})
		}

		response = append(response, models.TransactionResponse{
			ID:            transaction.ID,
			TotalPrice:    transaction.TotalPrice,
			CreatedAt:     transaction.CreatedAt,
			PaymentMethod: transaction.PaymentMethod,
			PaymentStatus: transaction.PaymentStatus,
			PaymentAt:     transaction.PaymentAt,
			Details:       detailResponse,
		})
	}

	if len(response) == 0 {
		response = []models.TransactionResponse{}
	}

	return response, http.StatusOK, nil
}

// add transaction from cart and remove cart after that, if cart is empty then return error
func (r *transactionRepository) AddTransaction(userId int) (models.Transaction, int, error) {
	var carts []models.Cart
	err := r.DB.Where("user_id = ?", userId).Find(&carts).Error
	if err != nil {
		return models.Transaction{}, http.StatusInternalServerError, err
	}

	if len(carts) == 0 {
		return models.Transaction{}, http.StatusBadRequest, errors.New("Cart is empty")
	}

	var totalPrice float64
	for _, cart := range carts {
		var tree models.Tree
		err := r.DB.Where("id = ?", cart.TreeID).First(&tree).Error
		if err != nil {
			return models.Transaction{}, http.StatusInternalServerError, err
		}

		totalPrice += tree.Price * float64(cart.Quantity)
	}

	transaction := models.Transaction{
		UserID:        userId,
		TotalPrice:    totalPrice,
		CreatedAt:     time.Now(),
		PaymentStatus: "pending",
	}

	err = r.DB.Create(&transaction).Error
	if err != nil {
		return models.Transaction{}, http.StatusInternalServerError, err
	}

	for _, cart := range carts {
		var tree models.Tree
		err := r.DB.Where("id = ?", cart.TreeID).First(&tree).Error
		if err != nil {
			return models.Transaction{}, http.StatusInternalServerError, err
		}

		detail := models.TransactionDetail{
			TransactionID: transaction.ID,
			TreeID:        tree.ID,
			Quantity:      cart.Quantity,
			Price:         tree.Price,
			TotalPrice:    tree.Price * float64(cart.Quantity),
		}

		err = r.DB.Create(&detail).Error
		if err != nil {
			return models.Transaction{}, http.StatusInternalServerError, err
		}

		err = r.DB.Delete(&cart).Error
		if err != nil {
			return models.Transaction{}, http.StatusInternalServerError, err
		}
	}

	return transaction, http.StatusOK, nil
}
