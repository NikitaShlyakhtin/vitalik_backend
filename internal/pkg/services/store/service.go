package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"time"
	"vitalik_backend/.gen/vitalik/public/table"
	"vitalik_backend/internal/dependencies"
	store_types "vitalik_backend/internal/pkg/services/store/types"
	"vitalik_backend/internal/pkg/types"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

var _ dependencies.IStore = (*Store)(nil)

func (s *Store) CreateWallet(ctx context.Context, wallet types.Wallet) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db.Begin failed: %w", err)
	}
	defer tx.Rollback(ctx)

	sql, args := table.Wallets.
		INSERT(table.Wallets.AllColumns).
		MODEL(store_types.MapToWalletStore(&wallet)).
		Sql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}

	transaction := types.Transaction{
		ID: uuid.NewString(),
		ReceiverRequisites: types.Requisites{
			Address: wallet.Requisites.Address,
			UserID:  wallet.Requisites.UserID,
		},
		Amount:    wallet.Balance,
		Currency:  wallet.Currency,
		Purpose:   null.StringFrom("Opening new wallet"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	sql, args = table.Transactions.
		INSERT(table.Transactions.AllColumns).
		MODEL(store_types.MapToTransactionStore(&transaction)).
		Sql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) ListWallets(ctx context.Context, args store_types.ListWalletsArgs) ([]*types.Wallet, error) {
	predicates := make([]postgres.BoolExpression, 0)

	if len(args.UserIDsIn) > 0 {
		predicates = append(predicates, table.Wallets.UserID.IN(
			lo.Map(args.UserIDsIn, func(userID string, _ int) postgres.Expression {
				return postgres.String(userID)
			})...,
		))
	}
	if len(args.AddresssIn) > 0 {
		predicates = append(predicates, table.Wallets.Address.IN(
			lo.Map(args.AddresssIn, func(address string, _ int) postgres.Expression {
				return postgres.String(address)
			})...,
		))
	}

	query := table.Wallets.
		SELECT(table.Wallets.AllColumns)

	if len(predicates) > 0 {
		query = query.
			WHERE(postgres.AND(predicates...))
	}

	sql, queryArgs := query.Sql()

	wallets := []store_types.Wallet{}
	if err := pgxscan.Select(ctx, s.db, &wallets, sql, queryArgs...); err != nil {
		return nil, fmt.Errorf("pgxscan.ScanAll failed: %w", err)
	}

	return lo.Map(wallets, func(wallet store_types.Wallet, _ int) *types.Wallet {
		return store_types.MapToWallet(wallet)
	}), nil
}

func (s *Store) Deposit(ctx context.Context, args store_types.DepositArgs) (*types.Transaction, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.Begin failed: %w", err)
	}
	defer tx.Rollback(ctx)

	sql, queryArgs := table.Wallets.
		UPDATE(table.Wallets.Balance, table.Wallets.UpdatedAt).
		SET(
			table.Wallets.Balance.ADD(postgres.Float(args.Amount)),
			postgres.TimestampzT(args.UpdatedAt),
		).
		WHERE(table.Wallets.Address.EQ(postgres.String(args.Address))).
		RETURNING(table.Wallets.AllColumns).
		Sql()

	fmt.Println(sql, args)

	rows, err := s.db.Query(ctx, sql, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("s.db.Query failed: %w", err)
	}
	defer rows.Close()

	wallets := []store_types.Wallet{}
	if err = pgxscan.ScanAll(&wallets, rows); err != nil {
		return nil, fmt.Errorf("pgxscan.ScanAll failed: %w", err)
	}

	if len(wallets) == 0 {
		return nil, ErrNotFound
	}

	wallet := wallets[0]

	transaction := &types.Transaction{
		ID: uuid.NewString(),
		ReceiverRequisites: types.Requisites{
			Address: args.Address,
			UserID:  wallet.UserID,
		},
		Amount:    args.Amount,
		Currency:  args.Currency,
		Purpose:   null.StringFrom("Deposit from external wallet"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	sql, queryArgs = table.Transactions.
		INSERT(table.Transactions.AllColumns).
		MODEL(store_types.MapToTransactionStore(transaction)).
		Sql()

	_, err = tx.Exec(ctx, sql, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("tx.Exec failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	return transaction, nil
}

func (s *Store) ListTransactions(ctx context.Context, args store_types.ListTransactionsArgs) ([]*types.Transaction, error) {
	predicates := make([]postgres.BoolExpression, 0)

	if len(args.AddresssIn) > 0 {
		predicates = append(predicates,
			postgres.OR(
				table.Transactions.SenderAddress.IN(
					lo.Map(args.AddresssIn, func(address string, _ int) postgres.Expression {
						return postgres.String(address)
					})...,
				),
				table.Transactions.ReceiverAddress.IN(
					lo.Map(args.AddresssIn, func(address string, _ int) postgres.Expression {
						return postgres.String(address)
					})...,
				),
			),
		)
	}

	if len(args.UserIDsIn) > 0 {
		predicates = append(predicates,
			postgres.OR(
				table.Transactions.SenderUserID.IN(
					lo.Map(args.UserIDsIn, func(userID string, _ int) postgres.Expression {
						return postgres.String(userID)
					})...,
				),
				table.Transactions.ReceiverUserID.IN(
					lo.Map(args.UserIDsIn, func(userID string, _ int) postgres.Expression {
						return postgres.String(userID)
					})...,
				),
			),
		)
	}

	query := table.Transactions.
		SELECT(table.Transactions.AllColumns)

	if len(predicates) > 0 {
		query = query.
			WHERE(postgres.AND(predicates...))
	}

	sql, queryArgs := query.Sql()

	transactions := []store_types.Transaction{}
	if err := pgxscan.Select(ctx, s.db, &transactions, sql, queryArgs...); err != nil {
		return nil, fmt.Errorf("pgxscan.ScanAll failed: %w", err)
	}

	return lo.Map(transactions, func(transaction store_types.Transaction, _ int) *types.Transaction {
		return store_types.MapToTransaction(&transaction)
	}), nil
}

func (s *Store) Transfer(ctx context.Context, args store_types.TransferArgs) (*types.Transaction, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.Begin failed: %w", err)
	}
	defer tx.Rollback(ctx)

	sql, queryArgs := table.Wallets.
		UPDATE(table.Wallets.Balance, table.Wallets.UpdatedAt).
		SET(
			table.Wallets.Balance.SUB(postgres.Float(args.Amount)),
			postgres.TimestampzT(time.Now()),
		).
		WHERE(
			postgres.AND(
				table.Wallets.Address.EQ(postgres.String(args.FromAddress)),
				table.Wallets.Balance.GT_EQ(postgres.Float(args.Amount)),
				table.Wallets.Currency.EQ(postgres.String(string(args.Currency))),
			),
		).
		RETURNING(table.Wallets.AllColumns).
		Sql()

	rows, err := s.db.Query(ctx, sql, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("s.db.Query failed: %w", err)
	}
	defer rows.Close()

	wallets := []store_types.Wallet{}
	if err = pgxscan.ScanAll(&wallets, rows); err != nil {
		return nil, fmt.Errorf("pgxscan.ScanAll failed: %w", err)
	}

	if len(wallets) == 0 {
		return nil, ErrNotFound
	}

	senderWallet := wallets[0]

	sql, queryArgs = table.Wallets.
		UPDATE(table.Wallets.Balance, table.Wallets.UpdatedAt).
		SET(
			table.Wallets.Balance.ADD(postgres.Float(args.Amount)),
			postgres.TimestampzT(time.Now()),
		).
		WHERE(
			postgres.AND(
				table.Wallets.Address.EQ(postgres.String(args.ToAddress)),
				table.Wallets.Currency.EQ(postgres.String(string(args.Currency))),
			),
		).
		RETURNING(table.Wallets.AllColumns).
		Sql()

	rows, err = s.db.Query(ctx, sql, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("s.db.Query failed: %w", err)
	}
	defer rows.Close()

	wallets = []store_types.Wallet{}
	if err = pgxscan.ScanAll(&wallets, rows); err != nil {
		return nil, fmt.Errorf("pgxscan.ScanAll failed: %w", err)
	}

	if len(wallets) == 0 {
		return nil, ErrNotFound
	}

	receiverWallet := wallets[0]

	transaction := &types.Transaction{
		ID: uuid.NewString(),
		SenderRequisites: types.Requisites{
			Address: senderWallet.Address,
			UserID:  senderWallet.UserID,
		},
		ReceiverRequisites: types.Requisites{
			Address: receiverWallet.Address,
			UserID:  receiverWallet.UserID,
		},
		Amount:    args.Amount,
		Currency:  args.Currency,
		Purpose:   args.Purpose,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	sql, queryArgs = table.Transactions.
		INSERT(table.Transactions.AllColumns).
		MODEL(store_types.MapToTransactionStore(transaction)).
		Sql()

	if _, err := tx.Exec(ctx, sql, queryArgs...); err != nil {
		return nil, fmt.Errorf("tx.Exec failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	return transaction, nil
}

func (s *Store) SaveOrder(ctx context.Context, order types.Order) error {
	sql, queryArgs := table.Orders.
		INSERT(table.Orders.AllColumns).
		MODEL(store_types.MapToOrderStore(&order)).
		Sql()

	if _, err := s.db.Exec(ctx, sql, queryArgs...); err != nil {
		return fmt.Errorf("s.db.Exec failed: %w", err)
	}

	return nil
}
