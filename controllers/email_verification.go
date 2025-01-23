package controllers

import (
	"carbon-api/repositories"
	"carbon-api/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EmailVerificationController struct {
	UserRepository repositories.UserRepository
}

func NewEmailVerificationController(userRepository repositories.UserRepository) *EmailVerificationController {
	return &EmailVerificationController{userRepository}
}

func (ctrl *EmailVerificationController) HandleEmailVerification(c echo.Context) error {
	token := c.QueryParam("token")
	fmt.Print("execute heree", token)

	if token == "" {
		return c.Render(http.StatusBadRequest, "verify-email.html", map[string]interface{}{
			"Message":      "Missing token. Please check your email verification link.",
			"MessageClass": "error",
		})
	}

	claims, err := utils.VerifyToken(token)
	if err != nil {
		return c.Render(http.StatusBadRequest, "verify-email.html", map[string]interface{}{
			"Message":      "Invalid or expired token. Please try again.",
			"MessageClass": "error",
		})
	}
	email := claims["email"].(string)

	user, status, err := ctrl.UserRepository.GetUserByEmail(email)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	if user.IsEmailVerified {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User already verified"})
	}

	updatedUser, err := ctrl.UserRepository.UpdateEmailVerificationStatus(user.ID, true)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	isVerified := updatedUser.IsEmailVerified
	messageClass := "error"
	message := "Email verification failed. Please try again."
	if isVerified {
		messageClass = "success"
		message = "Your email has been successfully verified! 🎉"
	}

	return c.Render(http.StatusOK, "verify-email.html", map[string]interface{}{
		"Message":      message,
		"MessageClass": messageClass,
	})
}
