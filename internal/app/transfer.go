package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	store_types "vitalik_backend/internal/pkg/services/store/types"
	"vitalik_backend/internal/pkg/types"
)

type transferRequest struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`

	Amount   float64        `json:"amount"`
	Currency types.Currency `json:"currency"`
}

func (app *Application) Transfer() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req *transferRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if code, err := app.validateTransferRequest(ctx, req); err != nil {
			return c.JSON(code, map[string]string{
				"message": fmt.Sprintf("validateTransferRequest failed: %v", err),
			})
		}

		args := bindTransferArgs(req)

		tx, err := app.Store.Transfer(ctx, args)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("Store Transfer failed: %v", err),
			})
		}

		return c.JSON(http.StatusOK, tx)
	}
}

func (app *Application) validateTransferRequest(ctx context.Context, req *transferRequest) (int, error) {
	if req.FromAddress == "" {
		return http.StatusBadRequest, fmt.Errorf("from_address must be provided")
	}

	if req.ToAddress == "" {
		return http.StatusBadRequest, fmt.Errorf("to_address must be provided")
	}

	matched, err := regexp.MatchString("^[0-9a-fA-F]{64}$", req.FromAddress)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to validate from_address")
	} else if !matched {
		return http.StatusBadRequest, fmt.Errorf("invalid from_address format")
	}

	matched, err = regexp.MatchString("^[0-9a-fA-F]{64}$", req.ToAddress)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to validate to_address")
	} else if !matched {
		return http.StatusBadRequest, fmt.Errorf("invalid to_address format")
	}

	if req.Amount <= 0 {
		return http.StatusBadRequest, fmt.Errorf("invalid amount")
	}

	if req.Currency == "" {
		return http.StatusBadRequest, errors.New("currency must be provided")
	} else if !req.Currency.Validate() {
		return http.StatusBadRequest, fmt.Errorf("invalid currency: %s", req.Currency)
	}

	return http.StatusOK, nil
}

func bindTransferArgs(req *transferRequest) store_types.TransferArgs {
	return store_types.TransferArgs{
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Purpose:     null.StringFrom("Internal transfer between exchange wallets"),
	}
}
