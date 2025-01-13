// repositories/role.go
package repositories

import (
	"carbon-api/models"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type RoleRepository interface {
	GetAllRoles() ([]models.Role, int, error)
	GetRoleByID(id int) (models.Role, int, error)
	CreateRole(role models.RoleRequest) (models.Role, int, error)
	UpdateRole(id int, role models.RoleRequest) (models.Role, int, error)
	DeleteRole(id int) (int, error)
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(DB *gorm.DB) RoleRepository {
	return &roleRepository{DB}
}

func (repo *roleRepository) GetAllRoles() ([]models.Role, int, error) {
	var roles []models.Role
	result := repo.DB.Order("id asc").Find(&roles)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	if len(roles) == 0 {
		roles = []models.Role{}
	}

	return roles, http.StatusOK, nil
}

func (repo *roleRepository) GetRoleByID(id int) (models.Role, int, error) {
	var role models.Role
	result := repo.DB.First(&role, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || role.ID == 0 {
		return role, http.StatusNotFound, errors.New("Role not found")
	}

	return role, http.StatusOK, nil
}

func (repo *roleRepository) CreateRole(roleRequest models.RoleRequest) (models.Role, int, error) {
	newRole := models.Role{
		Name: roleRequest.Name,
	}

	result := repo.DB.Create(&newRole)
	if result.Error != nil {
		return models.Role{}, http.StatusInternalServerError, result.Error
	}

	return newRole, http.StatusCreated, nil
}

func (repo *roleRepository) UpdateRole(id int, roleRequest models.RoleRequest) (models.Role, int, error) {
	var existingRole models.Role
	result := repo.DB.First(&existingRole, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || existingRole.ID == 0 {
		return existingRole, http.StatusNotFound, errors.New("Role not found")
	}

	existingRole.Name = roleRequest.Name

	result = repo.DB.Save(&existingRole)
	if result.Error != nil {
		return models.Role{}, http.StatusInternalServerError, result.Error
	}

	return existingRole, http.StatusOK, nil
}

func (repo *roleRepository) DeleteRole(id int) (int, error) {
	result := repo.DB.Delete(&models.Role{}, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Role not found")
	}

	return http.StatusOK, nil
}
