package app

import (
	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	order_book_types "vitalik_backend/internal/pkg/services/order_book/types"
	"vitalik_backend/internal/pkg/types"
)

type createOrderRequest struct {
	Type types.OrderType `json:"type"`

	SellCurrency   types.Currency   `json:"sell_currency"`
	SellQuantity   null.Float       `json:"sell_quantity"`
	SellRequisites types.Requisites `json:"sell_requisites"`

	Price float64

	BuyCurrency   types.Currency   `json:"buy_currency"`
	BuyQuantity   null.Float       `json:"buy_quantity"`
	BuyRequisites types.Requisites `json:"buy_requisites"`
}

func (app *Application) CreateOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var req createOrderRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		args := order_book_types.CreateOrderArgs{
			Type:           req.Type,
			SellCurrency:   req.SellCurrency,
			SellQuantity:   req.SellQuantity,
			SellRequisites: req.SellRequisites,
			Price:          req.Price,
			BuyCurrency:    req.BuyCurrency,
			BuyQuantity:    req.BuyQuantity,
			BuyRequisites:  req.BuyRequisites,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		order, err := app.OrderBook.CreateOrder(ctx, args)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, order)
	}
}
