// controllers/user.go
package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"carbon-api/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserRepository repositories.UserRepository
}

func NewUserController(userRepository repositories.UserRepository) *UserController {
	return &UserController{userRepository}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "Register User"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/register [post]
func (ctrl *UserController) RegisterUser(c echo.Context) error {
	var request models.RegisterRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	if request.Name == "" || request.Email == "" || request.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required fields"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}
	fmt.Printf("RoleID set during registration: %d\n", request.RoleID)

	request.Password = string(hashedPassword)

	// Set default role_id to "customer" (2)
	request.RoleID = 2

	user, status, err := ctrl.UserRepository.CreateUser(request)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	response := models.RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
	}
	// Kirim email notifikasi
	subject := "Welcome to Carbon App!"
	body := fmt.Sprintf("Hi %s,\n\nWelcome to Carbon App! Your account has been successfully created.\n\nThank you!", response.Name)
	if err := utils.SendEmail(response.Email, subject, body); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to send email"})
	}

	return c.JSON(status, response)
}

// LoginUser godoc
// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login User"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func (ctrl *UserController) LoginUser(c echo.Context) error {
	var request models.LoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	if request.Email == "" || request.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required fields"})
	}

	user, status, err := ctrl.UserRepository.GetUserByEmail(request.Email)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	// Ambil nama role berdasarkan RoleID
	roleName, err := ctrl.UserRepository.GetRoleNameByID(user.RoleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve role name"})
	}

	response := map[string]interface{}{
		"role":  roleName,
		"token": tokenString,
	}

	return c.JSON(http.StatusOK, response)

}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve the profile information of the authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} models.UserProfileResponse
// @Failure 401 {object} map[string]string
// @Router /users/profile [get]
func (ctrl *UserController) GetProfile(c echo.Context) error {
	userClaims := c.Get("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	// Retrieve user details from the database
	user, status, err := ctrl.UserRepository.GetUserByID(userID)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Map user to UserProfileResponse
	response := models.UserProfileResponse{
		ID:        user.ID,
		RoleID:    user.RoleID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile fields excluding email
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.UpdateProfileRequest true "Update Profile"
// @Success 200 {object} models.UpdateProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/profile [put]
func (ctrl *UserController) UpdateProfile(c echo.Context) error {
	userClaims := c.Get("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	var request models.UpdateProfileRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	// Validasi Email tidak boleh diubah
	if request.Email != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email cannot be updated"})
	}

	// Hash password baru jika ada
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
		}
		request.Password = string(hashedPassword)
	}

	// Perbarui data pengguna
	user, status, err := ctrl.UserRepository.UpdateUserProfile(userID, request)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Format respons tanpa password
	response := models.UpdateProfileResponse{

		Name:      user.Name,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// LogoutUser godoc
// @Summary Logout user
// @Description Invalidate the user's authentication token
// @Tags Users
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/logout [post]
func (ctrl *UserController) LogoutUser(c echo.Context) error {
	// Ambil token dari Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing token"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token format"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update user password with the current and new password
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.UpdatePasswordRequest true "Update Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/update-password [put]
func (ctrl *UserController) UpdatePassword(c echo.Context) error {
	userClaims := c.Get("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	var request models.UpdatePasswordRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	// Validasi input password
	if request.CurrentPassword == "" || request.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Both current and new passwords are required"})
	}

	// Ambil data pengguna berdasarkan ID
	user, status, err := ctrl.UserRepository.GetUserByID(userID)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Periksa apakah password saat ini cocok
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Current password is incorrect"})
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash new password"})
	}

	// Perbarui password di database
	user.Password = string(hashedPassword)
	if err := ctrl.UserRepository.UpdatePassword(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password updated successfully"})
}
