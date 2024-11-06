package app

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Application) HealthCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	}
}
