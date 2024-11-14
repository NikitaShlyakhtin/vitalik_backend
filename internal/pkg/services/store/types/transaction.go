package store_types

import (
	"github.com/guregu/null/v5"
	"time"
	"vitalik_backend/internal/pkg/types"
)

type Transaction struct {
	ID string `db:"transactions.id"`

	SenderAddress string `db:"transactions.sender_address"`
	SenderUserID  string `db:"transactions.sender_user_id"`

	ReceiverAddress string `db:"transactions.receiver_address"`
	ReceiverUserID  string `db:"transactions.receiver_user_id"`

	Amount   float64 `db:"transactions.amount"`
	Currency string  `db:"transactions.currency"`

	Purpose null.String `db:"transactions.purpose"`

	CreatedAt time.Time `db:"transactions.created_at"`
	UpdatedAt time.Time `db:"transactions.updated_at"`
}

// MapToTransactionStore converts a Transaction struct to a store.Transaction struct.
func MapToTransactionStore(tx *types.Transaction) *Transaction {
	return &Transaction{
		ID:              tx.ID,
		SenderAddress:   tx.SenderRequisites.Address,
		SenderUserID:    tx.SenderRequisites.UserID,
		ReceiverAddress: tx.ReceiverRequisites.Address,
		ReceiverUserID:  tx.ReceiverRequisites.UserID,
		Amount:          tx.Amount,
		Currency:        string(tx.Currency),
		Purpose:         tx.Purpose,
		CreatedAt:       tx.CreatedAt,
		UpdatedAt:       tx.UpdatedAt,
	}
}

// MapToTransaction converts a Transaction struct to a Transaction struct.
func MapToTransaction(txStore *Transaction) *types.Transaction {
	return &types.Transaction{
		ID: txStore.ID,
		SenderRequisites: types.Requisites{
			Address: txStore.SenderAddress,
			UserID:  txStore.SenderUserID,
		},
		ReceiverRequisites: types.Requisites{
			Address: txStore.ReceiverAddress,
			UserID:  txStore.ReceiverUserID,
		},
		Amount:    txStore.Amount,
		Currency:  types.Currency(txStore.Currency),
		Purpose:   txStore.Purpose,
		CreatedAt: txStore.CreatedAt,
		UpdatedAt: txStore.UpdatedAt,
	}
}
