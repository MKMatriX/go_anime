package handlers

import (
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
	response := c.Response()
	defer c.SetResponse(response)

	if err := echo.BindBody(c, &request); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(request)
	if validationErrors != nil {
		return common.SendFailedValidateResponse(c, validationErrors)
	}

	return nil
}

func (h *Handler) bindIdParam(c *echo.Context, idParamRequest *requests.IdParamRequest) error {
	response := c.Response()
	defer c.SetResponse(response)

	err := echo.BindPathValues(c, idParamRequest)

	if err != nil {
		return common.SendBadRequestResponse(c, "Couldn't parse id: "+c.Param("ID"))
	}

	return nil
}
