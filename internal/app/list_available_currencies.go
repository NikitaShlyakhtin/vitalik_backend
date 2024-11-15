package app

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vitalik_backend/internal/pkg/types"
)

func (app *Application) ListAvailableCurrencies() echo.HandlerFunc {
	return func(c echo.Context) error {
		availableCurrencies := types.ListAvailableCurrencies()
		return c.JSON(http.StatusOK, availableCurrencies)
	}
}
