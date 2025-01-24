package repositories

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) GetAllCart(userID int) ([]models.GetCartsResponse, int, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.GetCartsResponse), args.Int(1), args.Error(2)
}

func (m *MockCartRepository) AddCart(cart models.AddCartRequest) (models.Cart, int, error) {
	args := m.Called(cart)
	return args.Get(0).(models.Cart), args.Int(1), args.Error(2)
}

func (m *MockCartRepository) DeleteCart(cartID int, userID int) (int, error) {
	args := m.Called(cartID, userID)
	return args.Int(0), args.Error(1)
}
