package store_types

import (
	"github.com/guregu/null/v5"
	"vitalik_backend/internal/pkg/types"
)

type ListWalletsArgs struct {
	AddresssIn []string
	UserIDsIn  []string
}

type DepositArgs struct {
	Address  string
	Currency types.Currency
	Amount   float64
}

type ListTransactionsArgs struct {
	AddresssIn []string
	UserIDsIn  []string
}

type TransferArgs struct {
	FromAddress string
	ToAddress   string
	Amount      float64
	Currency    types.Currency
	Purpose     null.String
}
