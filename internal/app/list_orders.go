package app

import (
	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	order_book_types "vitalik_backend/internal/pkg/services/order_book/types"
	"vitalik_backend/internal/pkg/types"
)

type listOrdersRequest struct {
	UserID      null.String                   `json:"user_id"`
	OrderType   null.Value[types.OrderType]   `json:"order_type"`
	OrderStatus null.Value[types.OrderStatus] `json:"order_status"`
	Currency    null.Value[types.Currency]    `json:"currency"`
}

func (app *Application) ListOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var req listOrdersRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		args := order_book_types.ListOrdersArgs{
			UserID:      req.UserID,
			OrderType:   req.OrderType,
			OrderStatus: req.OrderStatus,
			Currency:    req.Currency,
		}

		orders, err := app.OrderBook.ListOrders(ctx, args)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, orders)
	}
}
