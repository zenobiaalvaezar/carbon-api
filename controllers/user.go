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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		//TODO = ROLE NAme
		//"role_name": role.Name,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	response := models.LoginResponse{
		Token: tokenString,
	}

	return c.JSON(http.StatusOK, response)
}
func (ctrl *UserController) GetProfile(c echo.Context) error {
	// Ambil klaim dari context
	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
	}

	// Ambil user_id dari klaim
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID in token"})
	}

	// Konversi user_id ke int
	intUserID := int(userID)

	// Ambil data user dari database
	user, status, err := ctrl.UserRepository.GetUserByID(intUserID)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
func (ctrl *UserController) UpdateProfile(c echo.Context) error {
	userClaims := c.Get("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	var request models.UpdateProfileRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	// Hash password baru jika ada
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
		}
		request.Password = string(hashedPassword)
	}

	// Perbarui profil pengguna
	updatedUser, status, err := ctrl.UserRepository.UpdateUserProfile(userID, request)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedUser)
}

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

	// Simpan token ke blacklist (opsional)
	// Example: Redis, Database, etc. - Simpan tokenString ke penyimpanan

	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}
