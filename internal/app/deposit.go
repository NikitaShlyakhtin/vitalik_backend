package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"time"
	store_types "vitalik_backend/internal/pkg/services/store/types"
	"vitalik_backend/internal/pkg/types"
)

type depositRequest struct {
	Address  string         `json:"address"`
	Currency types.Currency `json:"currency"`
	Amount   float64        `json:"amount"`
}

func (app *Application) Deposit() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req *depositRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if code, err := app.validateDepositRequest(ctx, req); err != nil {
			return c.JSON(code, map[string]string{
				"message": fmt.Sprintf("validateDepositRequest failed: %v", err),
			})
		}

		args := bindDepositArgs(req)

		tx, err := app.Store.Deposit(ctx, args)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("Store Deposit failed: %v", err),
			})
		}

		return c.JSON(http.StatusOK, tx)
	}
}

func (app *Application) validateDepositRequest(ctx context.Context, req *depositRequest) (int, error) {
	matched, err := regexp.MatchString("^[0-9a-fA-F]{64}$", req.Address)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to validate address")
	} else if !matched {
		return http.StatusBadRequest, fmt.Errorf("invalid address format")
	}

	if req.Currency == "" {
		return http.StatusBadRequest, errors.New("currency must be provided")
	} else if !req.Currency.Validate() {
		return http.StatusBadRequest, fmt.Errorf("invalid currency: %s", req.Currency)
	}

	if req.Amount <= 0 {
		return http.StatusBadRequest, errors.New("amount must be greater than zero")
	}

	return http.StatusOK, nil
}

func bindDepositArgs(req *depositRequest) store_types.DepositArgs {
	return store_types.DepositArgs{
		Address:   req.Address,
		Currency:  req.Currency,
		Amount:    req.Amount,
		UpdatedAt: time.Now(),
	}
}
