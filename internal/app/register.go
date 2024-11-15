package app

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"vitalik_backend/internal/pkg/services/auth_service"
	auth_types "vitalik_backend/internal/pkg/services/auth_service/types"
)

type registerRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

func (app *Application) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req registerRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		err := app.AuthService.Register(ctx, auth_types.AuthArgs{
			UserID:   req.UserID,
			Password: req.Password,
		})
		if err != nil {
			if errors.Is(err, auth_service.ErrAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "user registered successfully"})
	}
}
