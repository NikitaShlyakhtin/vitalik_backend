package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"sort"
	"vitalik_backend/internal/pkg/services/store"
	store_types "vitalik_backend/internal/pkg/services/store/types"
)

type listWalletsRequest struct {
	AddresssIn []string `json:"addresss_in"`
	UserIDsIn  []string `json:"user_ids_in"`
}

func (app *Application) ListWallets() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req *listWalletsRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if req == nil {
			req = &listWalletsRequest{}
		}

		if code, err := app.validateListWalletsRequest(ctx, req); err != nil {
			return c.JSON(code, map[string]string{
				"message": fmt.Sprintf("validateCreateWalletRequest failed: %v", err),
			})
		}

		args := bindListWalletsArgs(req)

		wallets, err := app.Store.ListWallets(ctx, args)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "wallet not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("Store ListWallets failed: %v", err),
			})
		}

		sort.Slice(wallets, func(i, j int) bool {
			return wallets[i].UpdatedAt.After(wallets[j].UpdatedAt)
		})

		return c.JSON(http.StatusOK, wallets)
	}
}

func (app *Application) validateListWalletsRequest(ctx context.Context, req *listWalletsRequest) (int, error) {
	if req.AddresssIn != nil {
		for _, key := range req.AddresssIn {
			matched, err := regexp.MatchString("^[0-9a-fA-F]{64}$", key)
			if err != nil {
				return http.StatusInternalServerError, errors.New("failed to validate address")
			} else if !matched {
				return http.StatusBadRequest, fmt.Errorf("invalid address format")
			}
		}
	}

	if req.UserIDsIn != nil {
		for _, userID := range req.UserIDsIn {
			if len(userID) == 0 {
				return http.StatusBadRequest, fmt.Errorf("userID must be non-empty")
			}
		}
	}

	return http.StatusOK, nil
}

func bindListWalletsArgs(req *listWalletsRequest) store_types.ListWalletsArgs {
	return store_types.ListWalletsArgs{
		AddresssIn: req.AddresssIn,
		UserIDsIn:  req.UserIDsIn,
	}
}
