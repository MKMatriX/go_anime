package handlers

import (
	"errors"
	"go_anime/internal/requests"
	"go_anime/internal/services"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func (h *Handler) UserRegister(c *echo.Context) error {
	var request requests.UserRegisterRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"})
	}

	validationErrors := h.ValidateBodyRequest(request)
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	userService := services.NewUserService(h.db)

	_, err := userService.GetUserByLogin(request.Login)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "User already exists"})
	}

	user, err := userService.RegisterUser(&request)

	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UserLogin(c *echo.Context) error {
	var request requests.UserLoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"})
	}

	validationErrors := h.ValidateBodyRequest(request)
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	userService := services.NewUserService(h.db)
	user, err := userService.LoginUser(&request)

	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
