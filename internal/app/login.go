package app

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"vitalik_backend/internal/pkg/services/auth_service"
	auth_types "vitalik_backend/internal/pkg/services/auth_service/types"
)

type loginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

func (app *Application) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req loginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		resp, err := app.AuthService.Login(ctx, auth_types.AuthArgs{
			UserID:   req.UserID,
			Password: req.Password,
		})
		if err != nil {
			if errors.Is(err, auth_service.ErrUserNotFound) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, resp)
	}
}
