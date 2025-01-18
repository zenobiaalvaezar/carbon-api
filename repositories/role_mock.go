package repositories

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetAllRoles() ([]models.Role, int, error) {
	args := m.Called()
	return args.Get(0).([]models.Role), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) GetRoleByID(id int) (models.Role, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.Role), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) CreateRole(roleRequest models.RoleRequest) (models.Role, int, error) {
	args := m.Called(roleRequest)
	return args.Get(0).(models.Role), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) UpdateRole(id int, roleRequest models.RoleRequest) (models.Role, int, error) {
	args := m.Called(id, roleRequest)
	return args.Get(0).(models.Role), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) DeleteRole(id int) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}
