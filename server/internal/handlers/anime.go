package handlers

import (
	"go_anime/internal/common"
	"go_anime/internal/requests"
	"go_anime/internal/services"
	"strconv"

	"github.com/labstack/echo/v5"
)

func (h *Handler) AnimeList(c *echo.Context) error {
	service := services.NewAnimeService(h.db)
	animeList := service.List()
	return common.SendSuccessResponse(c, "Anime list", &animeList)
}

func (h *Handler) AnimeCreate(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	var request requests.AnimeCreateRequest
	if err := c.Bind(&request); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(request)
	if validationErrors != nil {
		return common.SendFailedValidateResponse(c, validationErrors)
	}

	anime, err := service.Create(&request)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Anime added", anime)
}

func (h *Handler) AnimeUpdate(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	strId := c.Param("ID")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return common.SendBadRequestResponse(c, "Couldn't parse id: "+strId)
	}

	var request requests.AnimeCreateRequest
	if err := c.Bind(&request); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(request)
	if validationErrors != nil {
		return common.SendFailedValidateResponse(c, validationErrors)
	}

	anime, err := service.Update(id, &request)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Anime updated", anime)
}

func (h *Handler) AnimeDelete(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	strId := c.Param("ID")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return common.SendBadRequestResponse(c, "Couldn't parse id: "+strId)
	}

	err = service.Delete(id)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Anime deleted", id)
}
