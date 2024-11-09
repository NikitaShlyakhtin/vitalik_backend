package wallet_service_types

import "vitalik_backend/internal/pkg/types"

type ValidateCurrencyResponse struct {
	Currency  types.Currency
	Violation string
}
