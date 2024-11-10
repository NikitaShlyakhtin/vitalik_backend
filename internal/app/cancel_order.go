package app

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"vitalik_backend/internal/pkg/types"
)

type removeOrderRequest struct {
	CurrencyPair types.CurrencyPair `json:"currency_pair,omitempty"`
	OrderID      string             `param:"id"`
}

func (app *Application) CancelOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var req removeOrderRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid order ID"})
		}

		if !req.CurrencyPair.Validate() {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid currency_pair"})
		}

		if req.OrderID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "order ID is required"})
		}

		orderID, err := uuid.Parse(req.OrderID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid order order ID"})
		}

		err = app.OrderBookManager.CancelOrder(ctx, req.CurrencyPair, orderID)
		if err != nil {
			if errors.Is(err, echo.ErrNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "order book not found"})
			}
			return c.JSON(
				http.StatusInternalServerError, map[string]string{
					"message": fmt.Sprintf("OrderBookManager.CancelOrder failed: %s", err.Error()),
				},
			)
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "order removed successfully"})
	}
}
