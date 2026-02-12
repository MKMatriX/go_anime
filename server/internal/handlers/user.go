package handlers

import (
	"errors"
	"go_anime/internal/common"
	"go_anime/internal/models"
	"go_anime/internal/requests"
	"go_anime/internal/services"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func (h *Handler) UserRegister(c *echo.Context) error {
	var request requests.UserRegisterRequest

	err := h.bindAndValidate(c, &request)
	if err != nil {
		return err
	}

	userService := services.NewUserService(h.db)

	_, err = userService.GetUserByLogin(request.Login)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendBadRequestResponse(c, "User already exists")
	}

	user, err := userService.RegisterUser(&request)

	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Registered", user)
}

func (h *Handler) UserLogin(c *echo.Context) error {
	var request requests.UserLoginRequest

	err := h.bindAndValidate(c, &request)
	if err != nil {
		return err
	}

	userService := services.NewUserService(h.db)
	user, err := userService.LoginUser(&request)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	accessToken, refreshToken, err := common.GenerateJWT(*user)
	if err != nil {
		return common.SendBadRequestResponse(c, "JWT "+err.Error())
	}
	go userService.SaveRefreshTokenHash(user, refreshToken)

	return common.SendSuccessResponse(c, "Login successful", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (h *Handler) UserJWTRefresh(c *echo.Context) error {
	var request requests.UserJWTRefreshRequest

	err := h.bindAndValidate(c, &request)
	if err != nil {
		return common.SendBadRequestResponse(c, "bind and validate error")
	}

	userService := services.NewUserService(h.db)
	user, err := userService.GetUserByRefreshToken(request.RefreshToken)

	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	accessToken, refreshToken, err := common.GenerateJWT(*user)

	if err != nil {
		return common.SendBadRequestResponse(c, "JWT "+err.Error())
	}
	go userService.SaveRefreshTokenHash(user, refreshToken)

	return common.SendSuccessResponse(c, "Refresh successful", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (h Handler) GetAuthenticatedUser(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, "Unexpected error")
	}

	return common.SendSuccessResponse(c, "ok", user)
}
