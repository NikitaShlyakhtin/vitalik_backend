package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	wallet_service_types "vitalik_backend/internal/pkg/services/wallet_service/types"
	"vitalik_backend/internal/pkg/types"
)

type createWalletRequest struct {
	UserID   string         `json:"user_id"`
	Currency types.Currency `json:"currency"`
}

func (app *Application) CreateWallet() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req *createWalletRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if code, err := app.validateCreateWalletRequest(ctx, req); err != nil {
			return c.JSON(code, map[string]string{
				"message": fmt.Sprintf("validateCreateWalletRequest failed: %v", err),
			})
		}

		args := bindCreateWalletArgs(req)

		wallet, err := app.WalletService.CreateWallet(ctx, args)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("WalletService CreateWallet failed: %v", err),
			})
		}

		return c.JSON(http.StatusOK, wallet)
	}
}

func (app *Application) validateCreateWalletRequest(ctx context.Context, req *createWalletRequest) (int, error) {
	if req.UserID == "" {
		return http.StatusBadRequest, errors.New("user_id must be provided")
	}

	if req.Currency == "" {
		return http.StatusBadRequest, errors.New("currency must be provided")
	} else if !req.Currency.Validate() {
		return http.StatusBadRequest, fmt.Errorf("invalid currency: %s", req.Currency)
	}

	return http.StatusOK, nil
}

func bindCreateWalletArgs(req *createWalletRequest) wallet_service_types.CreateWalletArgs {
	return wallet_service_types.CreateWalletArgs{
		UserID:   req.UserID,
		Currency: req.Currency,
	}
}
