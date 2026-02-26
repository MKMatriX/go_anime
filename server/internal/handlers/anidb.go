package handlers

import (
	"go_anime/internal/requests"
	services "go_anime/internal/services/main"
	"go_anime/internal/shared/common"

	"github.com/labstack/echo/v5"
)

func (h *Handler) AnimeAutocomplete(c *echo.Context) error {
	service := services.NewAnimeService(h.db)

	var request requests.AnimeAutocompleteRequest
	err := h.bindParam(c, &request)
	if err != nil {
		return nil
	}

	titles, err := service.AutocompleteSearch(request.Query)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Titles", titles)
}
