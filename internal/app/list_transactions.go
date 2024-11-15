package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	store_types "vitalik_backend/internal/pkg/services/store/types"
)

type ListTransactionsRequest struct {
	AddresssIn []string `json:"address_in"`
	UserIdIn   []string `json:"user_id_in"`
}

func (app *Application) ListTransactions() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req *ListTransactionsRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if code, err := app.validateListTransactionsRequest(ctx, req); err != nil {
			return c.JSON(code, map[string]string{
				"message": fmt.Sprintf("validateListTransactionsRequest failed: %v", err),
			})
		}

		args := bindListTransactionsArgs(req)

		transactions, err := app.Store.ListTransactions(ctx, args)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("Store ListTransactions failed: %v", err),
			})
		}

		return c.JSON(http.StatusOK, transactions)
	}
}

func (app *Application) validateListTransactionsRequest(ctx context.Context, req *ListTransactionsRequest) (int, error) {
	if len(req.AddresssIn) == 0 &&
		len(req.UserIdIn) == 0 {
		return http.StatusBadRequest, fmt.Errorf("address or user_id must be provided")
	}

	if len(req.AddresssIn) > 0 {
		for _, key := range req.AddresssIn {
			matched, err := regexp.MatchString("^[0-9a-fA-F]{64}$", key)
			if err != nil {
				return http.StatusInternalServerError, errors.New("failed to validate address")
			} else if !matched {
				return http.StatusBadRequest, fmt.Errorf("invalid address format")
			}
		}
	}

	if len(req.UserIdIn) > 0 {
		for _, userID := range req.UserIdIn {
			if len(userID) == 0 {
				return http.StatusBadRequest, fmt.Errorf("user_id must must be non-empty")
			}
		}
	}

	return http.StatusOK, nil
}

func bindListTransactionsArgs(req *ListTransactionsRequest) store_types.ListTransactionsArgs {
	return store_types.ListTransactionsArgs{
		AddresssIn: req.AddresssIn,
		UserIDsIn:  req.UserIdIn,
	}
}
