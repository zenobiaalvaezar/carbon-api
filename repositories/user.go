// repositories/user.go
package repositories

import (
	"carbon-api/models"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user models.RegisterRequest) (models.User, int, error)
	UpdateUserProfile(userID int, request models.UpdateProfileRequest) (models.User, int, error)
	GetUserByEmail(email string) (models.User, int, error)
	GetUserByID(id int) (models.User, int, error)
	UpdatePassword(user models.User) error
	GetRoleNameByID(roleID int) (string, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userRepository{DB}
}

func (repo *userRepository) CreateUser(userRequest models.RegisterRequest) (models.User, int, error) {

	newUser := models.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
		Phone:    userRequest.Phone,
		Address:  userRequest.Address,
		RoleID:   userRequest.RoleID,
	}
	fmt.Printf("User data before insert: %+v\n", newUser)

	result := repo.DB.Create(&newUser)
	if result.Error != nil {
		return models.User{}, http.StatusInternalServerError, result.Error
	}

	return newUser, http.StatusCreated, nil
}

func (repo *userRepository) UpdateUserProfile(userID int, request models.UpdateProfileRequest) (models.User, int, error) {
	var user models.User
	result := repo.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, http.StatusNotFound, errors.New("User not found")
		}
		return user, http.StatusInternalServerError, result.Error
	}

	// Perbarui data pengguna
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if request.Address != "" {
		user.Address = request.Address
	}

	// Simpan perubahan
	result = repo.DB.Save(&user)
	if result.Error != nil {
		return user, http.StatusInternalServerError, result.Error
	}

	return user, http.StatusOK, nil
}

func (repo *userRepository) GetUserByEmail(email string) (models.User, int, error) {
	var user models.User
	result := repo.DB.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, http.StatusNotFound, errors.New("User not found")
	}

	return user, http.StatusOK, nil
}

func (repo *userRepository) GetUserByID(id int) (models.User, int, error) {
	var user models.User
	result := repo.DB.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, http.StatusNotFound, errors.New("User not found")
	}

	return user, http.StatusOK, nil
}
func (repo *userRepository) UpdatePassword(user models.User) error {
	result := repo.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *userRepository) GetRoleNameByID(roleID int) (string, error) {
	var role models.Role
	result := repo.DB.First(&role, roleID)
	if result.Error != nil {
		return "", result.Error
	}
	return role.Name, nil
}
