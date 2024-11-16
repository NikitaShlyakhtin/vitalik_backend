package app

import (
	"errors"
	"fmt"
	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	order_book_types "vitalik_backend/internal/pkg/services/order_book/types"
	"vitalik_backend/internal/pkg/types"
)

type listOrdersRequest struct {
	CurrencyPair types.CurrencyPair `json:"currency_pair,omitempty"`

	UserIDIn    []string                      `json:"user_id_in,omitempty"`
	OrderType   null.Value[types.OrderType]   `json:"order_type,omitempty"`
	OrderStatus null.Value[types.OrderStatus] `json:"order_status,omitempty"`
}

func (app *Application) ListOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var req listOrdersRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if !req.CurrencyPair.Validate() {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid currency_pair"})
		}

		args := order_book_types.ListOrdersArgs{
			CurrencyPair: req.CurrencyPair,
			UserIDIn:     req.UserIDIn,
			OrderType:    req.OrderType,
			OrderStatus:  req.OrderStatus,
		}

		orders, err := app.OrderBookManager.ListOrders(ctx, args)
		if err != nil {
			if errors.Is(err, echo.ErrNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "order book not found"})
			}
			return c.JSON(
				http.StatusInternalServerError, map[string]string{
					"message": fmt.Sprintf("OrderBookManager.ListOrders failed: %s", err.Error()),
				},
			)
		}

		return c.JSON(http.StatusOK, orders)
	}
}
