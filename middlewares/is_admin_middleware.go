package middlewares

import (
	"invitified-go/repositories"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func IsAdmin(repo repositories.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDStr, ok := c.Get("userID").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
			}

			user, err := repo.FindByID(userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
			}
			role, err := repo.FindRoleByID(user.RoleID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Role not found"})
			}

			if role.Name != "ADMIN" {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "Only admins can perform this action"})
			}

			return next(c)
		}
	}
}
