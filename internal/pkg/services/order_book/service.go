package order_book

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"sort"
	"time"

	"github.com/samber/lo"

	"vitalik_backend/internal/dependencies"
	order_book_types "vitalik_backend/internal/pkg/services/order_book/types"
	store_types "vitalik_backend/internal/pkg/services/store/types"
	"vitalik_backend/internal/pkg/types"
)

type OrderBook struct {
	store dependencies.IStore

	sellOrders []*types.Order
	buyOrders  []*types.Order
}

func NewOrderBook(
	store dependencies.IStore,
) (*OrderBook, error) {
	if store == nil {
		return nil, errors.New("failed to initialize order book")
	}

	return &OrderBook{
		store: store,

		sellOrders: make([]*types.Order, 0),
		buyOrders:  make([]*types.Order, 0),
	}, nil
}

func (o *OrderBook) CreateOrder(ctx context.Context, args order_book_types.CreateOrderArgs) (*types.Order, error) {
	order, err := bindOrderFromArgs(args)
	if err != nil {
		return nil, fmt.Errorf("BindOrderFromArgs failed: %w", err)
	}

	if err = o.validateOrder(ctx, order); err != nil {
		return nil, err
	}

	if err = o.createOrder(order); err != nil {
		return nil, fmt.Errorf("createOrder failed: %w", err)
	}

	return order, nil
}

func (o *OrderBook) MatchOrders(ctx context.Context) error {
	i := 0
	j := 0
	for i < len(o.buyOrders) && j < len(o.sellOrders) {
		buyOrder := *o.buyOrders[i]
		sellOrder := *o.sellOrders[j]

		if buyOrder.Status != types.OrderOpen {
			i++
			continue
		} else if sellOrder.Status != types.OrderClosed {
			j++
			continue
		}

		if buyOrder.Price >= sellOrder.Price {
			sellQuantity := min(buyOrder.BuyQuantity.ValueOrZero(), sellOrder.SellQuantity.ValueOrZero())

			transferArgs := store_types.TransferArgs{
				FromAddress: sellOrder.SellRequisites.Address,
				ToAddress:   buyOrder.BuyRequisites.Address,
				Amount:      sellQuantity,
				Currency:    sellOrder.SellCurrency,
				Purpose:     null.StringFrom("Trade"),
			}

			_, err := o.store.Transfer(ctx, transferArgs)
			if err != nil {
				return fmt.Errorf("store.Transfer failed: %w", err)
			}

			buyQuantity := sellQuantity * sellOrder.Price

			transferArgs = store_types.TransferArgs{
				FromAddress: buyOrder.SellRequisites.Address,
				ToAddress:   sellOrder.BuyRequisites.Address,
				Amount:      buyQuantity,
				Currency:    sellOrder.BuyCurrency,
				Purpose:     null.StringFrom("Trade"),
			}

			_, err = o.store.Transfer(ctx, transferArgs)
			if err != nil {
				return fmt.Errorf("store.Transfer failed: %w", err)
			}

			buyOrder.BuyQuantity = null.FloatFrom(buyOrder.BuyQuantity.ValueOrZero() - sellQuantity)
			sellOrder.SellQuantity = null.FloatFrom(sellOrder.SellQuantity.ValueOrZero() - sellQuantity)

			if buyOrder.BuyQuantity.ValueOrZero() <= 0 {
				buyOrder.Status = types.OrderClosed
			}
			if sellOrder.SellQuantity.ValueOrZero() <= 0 {
				sellOrder.Status = types.OrderClosed
			}

			o.buyOrders[i] = &buyOrder
			o.sellOrders[j] = &sellOrder
		} else {
			break
		}
	}

	return nil
}

func (o *OrderBook) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	for i, buyOrder := range o.buyOrders {
		if buyOrder.ID == orderID {
			buyOrder.Status = types.OrderCancelled
			o.buyOrders[i] = buyOrder
			return nil
		}
	}

	for i, sellOrder := range o.sellOrders {
		if sellOrder.ID == orderID {
			sellOrder.Status = types.OrderCancelled
			o.sellOrders[i] = sellOrder
			return nil
		}
	}

	return nil
}

func (o *OrderBook) ListOrders(ctx context.Context, args order_book_types.ListOrdersArgs) ([]*types.Order, error) {
	allOrders := append([]*types.Order{}, o.sellOrders...)
	allOrders = append(allOrders, o.buyOrders...)

	return lo.Filter(allOrders, func(order *types.Order, _ int) bool {
		if args.OrderType.Valid &&
			args.OrderType.ValueOrZero() == order.Type {
			return true
		}

		if args.OrderStatus.Valid &&
			args.OrderStatus.ValueOrZero() == order.Status {
			return true
		}

		if args.Currency.Valid &&
			(args.Currency.ValueOrZero() == order.SellCurrency ||
				args.Currency.ValueOrZero() == order.BuyCurrency) {
			return true
		}

		if args.UserID.Valid &&
			(args.UserID.ValueOrZero() == order.SellRequisites.UserID ||
				args.UserID.ValueOrZero() == order.BuyRequisites.UserID) {
			return true
		}

		return false
	}), nil
}

func (o *OrderBook) validateOrder(ctx context.Context, order *types.Order) error {
	wallets, err := o.getOrderWallets(ctx, order)
	if err != nil {
		return fmt.Errorf("getOrderWallets failed: %w", err)
	}

	if err = o.validateSellWallet(order, wallets); err != nil {
		return fmt.Errorf("validateSellWallet failed: %w", err)
	}

	if err = o.validateBuyWallet(order, wallets); err != nil {
		return fmt.Errorf("validateBuyWallet failed: %w", err)
	}

	return nil
}

func (o *OrderBook) getOrderWallets(ctx context.Context, order *types.Order) ([]*types.Wallet, error) {
	listWalletsArgs := store_types.ListWalletsArgs{
		AddresssIn: []string{
			order.SellRequisites.Address,
			order.BuyRequisites.Address,
		},
		UserIDsIn: []string{
			order.SellRequisites.UserID,
			order.BuyRequisites.UserID,
		},
	}

	wallets, err := o.store.ListWallets(ctx, listWalletsArgs)
	if err != nil {
		return nil, fmt.Errorf("store.ListWallets failed: %w", err)
	}

	return wallets, nil
}

func (o *OrderBook) validateSellWallet(order *types.Order, wallets []*types.Wallet) error {
	sellRequisites, found := findWallet(wallets, order.SellRequisites.Address)
	if !found {
		return fmt.Errorf("failed to find sell wallet: %v", order.SellRequisites.Address)
	}
	if sellRequisites.Currency != order.SellCurrency {
		return fmt.Errorf(
			"currency mismatch for sell wallet: %v, want %s, got %s",
			sellRequisites.Requisites.Address,
			sellRequisites.Currency,
			order.SellCurrency,
		)
	}
	if sellRequisites.Balance < order.SellQuantity.ValueOrZero() {
		return fmt.Errorf("insufficient balance for sell wallet: %v", sellRequisites.Requisites.Address)
	}
	return nil
}

func (o *OrderBook) validateBuyWallet(order *types.Order, wallets []*types.Wallet) error {
	BuyRequisites, found := findWallet(wallets, order.BuyRequisites.Address)
	if !found {
		return fmt.Errorf("failed to find want currency wallet: %v", order.BuyRequisites.Address)
	}
	if BuyRequisites.Currency != order.BuyCurrency {
		return fmt.Errorf(
			"currency mismatch for want currency wallet: %v, want %s, got %s",
			BuyRequisites.Requisites.Address,
			BuyRequisites.Currency,
			order.BuyCurrency,
		)
	}
	return nil
}

func findWallet(wallets []*types.Wallet, Address string) (*types.Wallet, bool) {
	return lo.Find(wallets, func(wallet *types.Wallet) bool {
		return wallet.Requisites.Address == Address
	})
}

func (o *OrderBook) createOrder(order *types.Order) error {
	if order.Type == types.Sell {
		o.sellOrders = append(o.sellOrders, order)
		o.sortOrderBook()
		return nil
	} else if order.Type == types.Buy {
		o.buyOrders = append(o.buyOrders, order)
		o.sortOrderBook()
		return nil
	}

	return fmt.Errorf("invalid order type: %v", order.Type)
}

func (o *OrderBook) sortOrderBook() {
	sort.Slice(o.sellOrders, func(i, j int) bool {
		o1 := o.sellOrders[i]
		o2 := o.sellOrders[j]

		if o1.Price == o2.Price {
			return o1.CreatedAt.Before(o2.CreatedAt)
		}
		return o1.Price < o2.Price
	})

	sort.Slice(o.buyOrders, func(i, j int) bool {
		o1 := o.buyOrders[i]
		o2 := o.buyOrders[j]

		if o1.Price == o2.Price {
			return o1.CreatedAt.Before(o2.CreatedAt)
		}
		return o1.Price > o2.Price
	})
}

func bindOrderFromArgs(args order_book_types.CreateOrderArgs) (*types.Order, error) {
	return &types.Order{
		ID:             uuid.Must(uuid.NewV7()),
		Type:           args.Type,
		SellCurrency:   args.SellCurrency,
		SellQuantity:   args.SellQuantity,
		SellRequisites: args.SellRequisites,
		Price:          args.Price,
		BuyCurrency:    args.BuyCurrency,
		BuyQuantity:    args.BuyQuantity,
		BuyRequisites:  args.BuyRequisites,
		Status:         types.OrderOpen,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
