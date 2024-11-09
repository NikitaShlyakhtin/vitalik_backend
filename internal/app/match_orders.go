package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Application) MatchOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		err := app.OrderBook.MatchOrders(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("failed to match orders: %v", err),
			})
		}

		// Return a success response if matching was successful
		return c.JSON(http.StatusOK, map[string]string{"message": "Orders matched successfully"})
	}
}
