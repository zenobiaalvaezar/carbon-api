package repositories

import (
	"carbon-api/models"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetReportSummary(userId int) ([]models.ReportSummary, int, error)
	GetReportDetail(userId int) ([]models.ReportDetail, int, error)
}

type reportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(DB *gorm.DB) *reportRepository {
	return &reportRepository{DB}
}

func (r *reportRepository) GetReportSummary(userId int) ([]models.ReportSummary, int, error) {
	var user models.User
	r.DB.Where("id = ?", userId).First(&user)
	if user.ID == 0 {
		return nil, http.StatusNotFound, errors.New("User not found")
	}

	var carbonTree int
	err := r.DB.Raw("SELECT SUM(total_tree) FROM carbon_summaries WHERE user_id = ?", userId).Scan(&carbonTree).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var donationTree int
	err = r.DB.Raw("SELECT SUM(quantity) FROM transaction_details WHERE transaction_id IN (SELECT id FROM transactions WHERE user_id = ? AND payment_status = 'success')", userId).Scan(&donationTree).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	badgeStatus := "Netral"
	if carbonTree > donationTree {
		badgeStatus = "Perusak Alam"
	} else if donationTree > carbonTree {
		badgeStatus = "Pelindung Alam"
	}

	response := []models.ReportSummary{
		{
			UserID:       userId,
			UserName:     user.Name,
			UserEmail:    user.Email,
			CarbonTree:   carbonTree,
			DonationTree: donationTree,
			BadgeStatus:  badgeStatus,
		},
	}

	return response, http.StatusOK, nil
}

func (r *reportRepository) GetReportDetail(userId int) ([]models.ReportDetail, int, error) {
	var user models.User
	r.DB.Where("id = ?", userId).First(&user)
	if user.ID == 0 {
		return nil, http.StatusNotFound, errors.New("User not found")
	}

	var transactions []models.Transaction
	r.DB.Where("user_id = ?", userId).Find(&transactions)
	if len(transactions) == 0 {
		return nil, http.StatusNotFound, errors.New("Transaction not found")
	}

	var response []models.ReportDetail
	for _, transaction := range transactions {
		var details []models.TransactionDetail
		r.DB.Where("transaction_id = ?", transaction.ID).Find(&details)

		var totalTree int
		for _, detail := range details {
			totalTree += detail.Quantity
		}

		response = append(response, models.ReportDetail{
			UserID:          userId,
			UserName:        user.Name,
			UserEmail:       user.Email,
			TransactionID:   transaction.ID,
			TransactionDate: transaction.CreatedAt,
			TotalTree:       totalTree,
			TotalPrice:      transaction.TotalPrice,
			PaymentMethod:   transaction.PaymentMethod,
			PaymentStatus:   transaction.PaymentStatus,
			PaymentAt:       transaction.PaymentAt,
		})
	}

	return response, http.StatusOK, nil
}
