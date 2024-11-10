package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/samber/lo"
	"time"
	"vitalik_backend/internal/dependencies"
	store_types "vitalik_backend/internal/pkg/services/store/types"
	"vitalik_backend/internal/pkg/types"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store struct {
	wallets      map[string]*types.Wallet
	transactions []*types.Transaction
	orders       []*types.Order
}

func NewStore() *Store {
	return &Store{
		wallets: make(map[string]*types.Wallet),
	}
}

var _ dependencies.IStore = (*Store)(nil)

func (s *Store) CreateWallet(ctx context.Context, wallet types.Wallet) error {
	s.wallets[wallet.Requisites.Address] = &wallet

	tx := &types.Transaction{
		ID: uuid.NewString(),
		ReceiverRequisites: types.Requisites{
			Address: wallet.Requisites.Address,
			UserID:  wallet.Requisites.UserID,
		},
		Amount:    wallet.Balance,
		Currency:  wallet.Currency,
		Purpose:   null.StringFrom("wallet creation"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.transactions = append(s.transactions, tx)

	return nil
}

func (s *Store) ListWallets(ctx context.Context, args store_types.ListWalletsArgs) ([]*types.Wallet, error) {
	return lo.Filter(lo.Values(s.wallets), func(wallet *types.Wallet, _ int) bool {
		Address := wallet.Requisites.Address
		userID := wallet.Requisites.UserID

		if len(args.UserIDsIn) > 0 && !lo.Contains(args.UserIDsIn, userID) {
			return false
		}

		if len(args.AddresssIn) > 0 && !lo.Contains(args.AddresssIn, Address) {
			return false
		}

		return true
	}), nil
}

func (s *Store) Deposit(ctx context.Context, args store_types.DepositArgs) (*types.Transaction, error) {
	wallet, ok := s.wallets[args.Address]
	if !ok {
		return nil, ErrNotFound
	} else if wallet.Currency != args.Currency {
		return nil, fmt.Errorf(
			"currency mismatch for wallet: %v, want %s, got %s",
			wallet.Requisites.Address,
			wallet.Currency,
			args.Currency,
		)
	}

	wallet.Balance += args.Amount
	wallet.UpdatedAt = time.Now()

	s.wallets[args.Address] = wallet

	tx := &types.Transaction{
		ID: uuid.NewString(),
		ReceiverRequisites: types.Requisites{
			Address: wallet.Requisites.Address,
			UserID:  wallet.Requisites.UserID,
		},
		Amount:    args.Amount,
		Currency:  args.Currency,
		Purpose:   null.StringFrom("deposit from external wallet"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.transactions = append(s.transactions, tx)

	return tx, nil
}

func (s *Store) ListTransactions(ctx context.Context, args store_types.ListTransactionsArgs) ([]*types.Transaction, error) {
	return lo.Filter(s.transactions, func(tx *types.Transaction, _ int) bool {
		senderAddress := tx.SenderRequisites.Address
		senderUserID := tx.SenderRequisites.UserID

		receiverAddress := tx.ReceiverRequisites.Address
		receiverUserID := tx.ReceiverRequisites.UserID

		return lo.Contains(args.AddresssIn, senderAddress) ||
			lo.Contains(args.AddresssIn, receiverAddress) ||
			lo.Contains(args.UserIDsIn, senderUserID) ||
			lo.Contains(args.UserIDsIn, receiverUserID)
	}), nil
}

func (s *Store) Transfer(ctx context.Context, args store_types.TransferArgs) (*types.Transaction, error) {
	from, ok := s.wallets[args.FromAddress]
	if !ok {
		return nil,
			fmt.Errorf("from wallet does not exist: %w", ErrNotFound)
	} else if from.Balance < args.Amount {
		return nil,
			fmt.Errorf("insufficient funds for wallet %v: %d < %d", from.Requisites.Address, from.Balance, args.Amount)
	} else if from.Currency != args.Currency {
		return nil, fmt.Errorf(
			"currency mismatch for sender wallet: %v, want %s, got %s",
			from.Requisites.Address,
			from.Currency,
			args.Currency,
		)
	}

	to, ok := s.wallets[args.ToAddress]
	if !ok {
		return nil,
			fmt.Errorf("to wallet does not exist: %w", ErrNotFound)
	} else if to.Currency != args.Currency {
		return nil, fmt.Errorf(
			"currency mismatch for receiver wallet: %v, want %s, got %s",
			to.Requisites.Address,
			to.Currency,
			args.Currency,
		)
	}

	from.Balance -= args.Amount
	from.UpdatedAt = time.Now()

	to.Balance += args.Amount
	to.UpdatedAt = time.Now()

	s.wallets[args.FromAddress] = from
	s.wallets[args.ToAddress] = to

	tx := &types.Transaction{
		ID: uuid.NewString(),
		ReceiverRequisites: types.Requisites{
			Address: from.Requisites.Address,
			UserID:  from.Requisites.UserID,
		},
		SenderRequisites: types.Requisites{
			Address: to.Requisites.Address,
			UserID:  to.Requisites.UserID,
		},
		Amount:    args.Amount,
		Currency:  args.Currency,
		Purpose:   args.Purpose,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.transactions = append(s.transactions, tx)

	return tx, nil
}

func (s *Store) SaveOrder(ctx context.Context, order types.Order) error {
	s.orders = append(s.orders, &order)

	return nil
}
