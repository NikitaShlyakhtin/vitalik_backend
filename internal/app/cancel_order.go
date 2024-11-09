package app

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type removeOrderRequest struct {
	OrderID string `param:"id"`
}

func (app *Application) CancelOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var req removeOrderRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid order ID"})
		}

		if req.OrderID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "order ID is required"})
		}

		orderID, err := uuid.Parse(req.OrderID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid order order ID"})
		}

		err = app.OrderBook.CancelOrder(ctx, orderID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "order removed successfully"})
	}
}
