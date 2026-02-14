package handlers

import (
	"go_anime/internal/common"
	"go_anime/internal/requests"
	"go_anime/internal/services"

	"github.com/labstack/echo/v5"
)

func (h *Handler) AnimeList(c *echo.Context) error {
	service := services.NewAnimeService(h.db)
	animeList := service.List()
	return common.SendSuccessResponse(c, "Anime list", &animeList)
}

func (h *Handler) AnimeItem(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	var request requests.IdParamRequest
	h.bindIdParam(c, &request)

	anime, err := service.GetById(request.ID)
	if err != nil {
		return common.SendNotFoundResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Anime found", &anime)
}

func (h *Handler) AnimeCreate(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	var request requests.AnimeCreateRequest
	err := h.bindAndValidate(c, &request)
	if err != nil {
		return err
	}

	anime, err := service.Create(&request)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Anime added", anime)
}

func (h *Handler) AnimeUpdate(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	var idRequest requests.IdParamRequest
	err := h.bindIdParam(c, &idRequest)
	if err != nil {
		return err
	}

	var request requests.AnimeCreateRequest
	err = h.bindAndValidate(c, &request)
	if err != nil {
		return err
	}

	anime, err := service.Update(idRequest.ID, &request)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Anime updated", anime)
}

func (h *Handler) AnimeDelete(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	var idRequest requests.IdParamRequest
	err := h.bindIdParam(c, &idRequest)
	if err != nil {
		return err
	}

	err = service.Delete(idRequest.ID)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Anime deleted", idRequest.ID)
}

func (h *Handler) AnimeEpisodes(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	var idRequest requests.IdParamRequest
	err := h.bindIdParam(c, &idRequest)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	anime, err := service.GetById(idRequest.ID)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	episodes, err := service.GetEpisodes(anime)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Episodes", episodes)
}
