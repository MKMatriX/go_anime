package middleware

import (
	"go_anime/internal/shared/common"
	"go_anime/internal/shared/models"
	"strings"

	"github.com/labstack/echo/v5"
)

func (am *AppMiddleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Add("Vary", "Authorisation")
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return common.SendUnauthorizedResponse(c, "Please provide a Bearer token")
		}

		authHeaderSplit := strings.Split(authHeader, " ")
		if len(authHeaderSplit) < 2 {
			return common.SendUnauthorizedResponse(c, "Please provide a Bearer token")
		}
		accessToken := authHeaderSplit[1]
		claims, err := common.ParseJWT(accessToken)
		if err != nil {
			return common.SendUnauthorizedResponse(c, err.Error())
		}

		if common.IsClaimsExpired(claims) {
			return common.SendUnauthorizedResponse(c, "Token has expired")
		}

		var user models.UserModel

		result := am.DB.First(&user, claims.ID)
		if result.Error != nil {
			return common.SendUnauthorizedResponse(c, "Invalid token")
		}

		c.Set("user", user)

		return next(c)
	}
}
