package middlewares

import (
	"invitified-go/repositories"
	"invitified-go/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(tokenRepo repositories.TokenRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing token"})
			}

			claims, err := utils.ValidateJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
			}

			userID := claims.UserID
			c.Set("userID", userID)

			return next(c)
		}
	}
}
