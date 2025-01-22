package repositories

import (
	"carbon-api/models"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetAllCart(userID int) ([]models.GetCartsResponse, int, error)
	AddCart(cart models.AddCartRequest) (models.Cart, int, error)
	DeleteCart(cartID int) (int, error)
}

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(DB *gorm.DB) *cartRepository {
	return &cartRepository{DB}
}

func (r *cartRepository) GetAllCart(userID int) ([]models.GetCartsResponse, int, error) {
	var carts []models.Cart
	err := r.DB.Where("user_id = ?", userID).Find(&carts).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var response []models.GetCartsResponse
	for _, cart := range carts {
		var tree models.Tree
		err := r.DB.Where("id = ?", cart.TreeID).First(&tree).Error
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		response = append(response, models.GetCartsResponse{
			ID:       cart.ID,
			TreeName: tree.Name,
			Price:    tree.Price,
			Quantity: cart.Quantity,
		})
	}

	if len(response) == 0 {
		response = []models.GetCartsResponse{}
	}

	return response, http.StatusOK, nil
}

func (r *cartRepository) AddCart(cart models.AddCartRequest) (models.Cart, int, error) {
	// check if tree exists
	var tree models.Tree
	r.DB.Where("id = ?", cart.TreeID).First(&tree)
	if tree.ID == 0 {
		return models.Cart{}, http.StatusNotFound, errors.New("Tree not found")
	}

	// check if tree stock enough
	if tree.Stock < cart.Quantity {
		return models.Cart{}, http.StatusBadRequest, errors.New("Tree stock not enough")
	}

	// check if cart already exists
	var existingCart models.Cart
	r.DB.Where("user_id = ? AND tree_id = ?", cart.UserID, cart.TreeID).First(&existingCart)
	if existingCart.ID != 0 {
		existingCart.Quantity += cart.Quantity
		err := r.DB.Save(&existingCart).Error
		if err != nil {
			return models.Cart{}, http.StatusInternalServerError, err
		}

		return models.Cart{}, http.StatusOK, nil
	}

	newCart := models.Cart{
		UserID:   cart.UserID,
		TreeID:   cart.TreeID,
		Quantity: cart.Quantity,
	}
	err := r.DB.Create(&newCart).Error
	if err != nil {
		return models.Cart{}, http.StatusInternalServerError, err
	}

	// update tree stock
	tree.Stock -= cart.Quantity
	err = r.DB.Save(&tree).Error
	if err != nil {
		return models.Cart{}, http.StatusInternalServerError, err
	}

	return newCart, http.StatusOK, nil
}

func (r *cartRepository) DeleteCart(cartID int) (int, error) {
	// check if cart exists
	var existingCart models.Cart
	r.DB.Where("id = ?", cartID).First(&existingCart)
	if existingCart.ID == 0 {
		return http.StatusNotFound, errors.New("Cart not found")
	}

	err := r.DB.Delete(&models.Cart{}, cartID).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
