package handlers

import (
	"errors"
	"go_anime/internal/common"
	"go_anime/internal/requests"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) bindAndValidate(c *echo.Context, request interface{}) error {
	if err := echo.BindBody(c, request); err != nil {
		common.SendBadRequestResponse(c, err.Error())
		return errors.New("Bind error")
	}

	validationErrors := h.ValidateBodyRequest(request)
	if validationErrors != nil {
		common.SendFailedValidateResponse(c, validationErrors)
		return errors.New("Validations errors")
	}

	return nil
}

func (h *Handler) bindIdParam(c *echo.Context, idParamRequest *requests.IdParamRequest) error {
	err := echo.BindPathValues(c, idParamRequest)

	if err != nil {
		common.SendBadRequestResponse(c, "Couldn't parse id: "+c.Param("ID"))
		return errors.New("Binding Id param error")
	}

	return nil
}
