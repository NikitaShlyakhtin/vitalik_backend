package order_book_manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"sort"
	"time"
	"vitalik_backend/internal/dependencies"
	"vitalik_backend/internal/pkg/services/order_book"
	order_book_types "vitalik_backend/internal/pkg/services/order_book/types"
	store_types "vitalik_backend/internal/pkg/services/store/types"
	"vitalik_backend/internal/pkg/types"
)

// OrderBookManager manages multiple order books in memory.
type OrderBookManager struct {
	store dependencies.IStore

	orderBooks map[string]*order_book.OrderBook
}

func NewOrderBookManager(store dependencies.IStore) (*OrderBookManager, error) {
	if store == nil {
		return nil, errors.New("failed to initialize order_book_manager")
	}

	return &OrderBookManager{
		store:      store,
		orderBooks: make(map[string]*order_book.OrderBook),
	}, nil
}

var _ dependencies.IOrderBookManager = (*OrderBookManager)(nil)

func (m *OrderBookManager) ListAvailableCurrencyPairs(ctx context.Context) ([]types.CurrencyPair, error) {
	pairs := lo.Map(lo.Values(m.orderBooks), func(book *order_book.OrderBook, _ int) types.CurrencyPair {
		return book.CurrencyPair
	})

	return pairs, nil
}

func (m *OrderBookManager) CreateOrder(ctx context.Context, args order_book_types.CreateOrderArgs) (*types.Order, error) {
	orderBook, err := m.getOrCreateOrderBook(ctx, types.CurrencyPair{
		Currency1: args.BuyCurrency,
		Currency2: args.SellCurrency,
	})
	if err != nil {
		return nil, fmt.Errorf("getOrCreateOrderBook failed: %w", err)
	}

	order, err := bindOrderFromArgs(args)
	if err != nil {
		return nil, fmt.Errorf("BindOrderFromArgs failed: %w", err)
	}

	if err = m.validateOrder(ctx, order); err != nil {
		return nil, err
	}

	if err = m.createOrder(orderBook, order); err != nil {
		return nil, fmt.Errorf("createOrder failed: %w", err)
	}

	return order, nil
}

func (m *OrderBookManager) CancelOrder(
	ctx context.Context,
	currencyPair types.CurrencyPair,
	orderID uuid.UUID,
) error {
	orderBook, err := m.getOrCreateOrderBook(ctx, currencyPair)
	if err != nil {
		return fmt.Errorf("getOrCreateOrderBook failed: %w", err)
	}

	for i, buyOrder := range orderBook.BuyOrders {
		if buyOrder.ID == orderID {
			buyOrder.Status = types.OrderCancelled
			orderBook.BuyOrders[i] = buyOrder
			return nil
		}
	}

	for i, sellOrder := range orderBook.SellOrders {
		if sellOrder.ID == orderID {
			sellOrder.Status = types.OrderCancelled
			orderBook.SellOrders[i] = sellOrder
			return nil
		}
	}

	return echo.ErrNotFound
}

func (m *OrderBookManager) ListOrders(ctx context.Context, args order_book_types.ListOrdersArgs) ([]*types.Order, error) {
	orderBook, err := m.getOrCreateOrderBook(ctx, args.CurrencyPair)
	if err != nil {
		return nil, fmt.Errorf("getOrderBook failed: %w", err)
	}

	allOrders := append([]*types.Order{}, orderBook.SellOrders...)
	allOrders = append(allOrders, orderBook.BuyOrders...)

	return lo.Filter(allOrders, func(order *types.Order, _ int) bool {
		if len(args.UserIDIn) > 0 &&
			!lo.Contains(args.UserIDIn, order.SellRequisites.UserID) &&
			!lo.Contains(args.UserIDIn, order.BuyRequisites.UserID) {
			return false
		}

		if args.OrderType.Valid &&
			args.OrderType.ValueOrZero() != order.Type {
			return false
		}

		if args.OrderStatus.Valid &&
			args.OrderStatus.ValueOrZero() != order.Status {
			return false
		}

		return true
	}), nil
}

func (m *OrderBookManager) MatchOrders(ctx context.Context) error {
	for key, book := range m.orderBooks {
		if err := m.matchOrders(ctx, book); err != nil {
			return fmt.Errorf("MatchOrders failed: %w", err)
		}
		m.orderBooks[key] = book
	}

	return nil
}

func (m *OrderBookManager) validateOrder(ctx context.Context, order *types.Order) error {
	wallets, err := m.getOrderWallets(ctx, order)
	if err != nil {
		return fmt.Errorf("getOrderWallets failed: %w", err)
	}

	if err = m.validateSellWallet(order, wallets); err != nil {
		return fmt.Errorf("validateSellWallet failed: %w", err)
	}

	if err = m.validateBuyWallet(order, wallets); err != nil {
		return fmt.Errorf("validateBuyWallet failed: %w", err)
	}

	return nil
}

func (m *OrderBookManager) getOrderWallets(ctx context.Context, order *types.Order) ([]*types.Wallet, error) {
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

	wallets, err := m.store.ListWallets(ctx, listWalletsArgs)
	if err != nil {
		return nil, fmt.Errorf("storefx.ListWallets failed: %w", err)
	}

	return wallets, nil
}

func (m *OrderBookManager) validateSellWallet(order *types.Order, wallets []*types.Wallet) error {
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

func (m *OrderBookManager) validateBuyWallet(order *types.Order, wallets []*types.Wallet) error {
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

func (m *OrderBookManager) createOrder(orderBook *order_book.OrderBook, order *types.Order) error {
	if order.Type == types.Sell {
		orderBook.SellOrders = append(orderBook.SellOrders, order)
		sortOrderBook(orderBook)
		return nil
	} else if order.Type == types.Buy {
		orderBook.BuyOrders = append(orderBook.BuyOrders, order)
		sortOrderBook(orderBook)
		return nil
	}

	return fmt.Errorf("invalid order type: %v", order.Type)
}

func sortOrderBook(orderBook *order_book.OrderBook) {
	sort.Slice(orderBook.SellOrders, func(i, j int) bool {
		o1 := orderBook.SellOrders[i]
		o2 := orderBook.SellOrders[j]

		if o1.Price == o2.Price {
			return o1.CreatedAt.Before(o2.CreatedAt)
		}
		return o1.Price < o2.Price
	})

	sort.Slice(orderBook.BuyOrders, func(i, j int) bool {
		o1 := orderBook.BuyOrders[i]
		o2 := orderBook.BuyOrders[j]

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

func (m *OrderBookManager) getOrCreateOrderBook(ctx context.Context, currencyPair types.CurrencyPair) (*order_book.OrderBook, error) {
	if lo.HasKey(m.orderBooks, currencyPair.String()) {
		return m.orderBooks[currencyPair.String()], nil
	} else if lo.HasKey(m.orderBooks, currencyPair.StringReverse()) {
		return m.orderBooks[currencyPair.StringReverse()], nil
	}

	orderBook, err := order_book.NewOrderBook(currencyPair)
	if err != nil {
		return nil, fmt.Errorf("NewOrderBook failed: %v", err)
	}

	m.orderBooks[currencyPair.String()] = orderBook

	return orderBook, nil
}

func (m *OrderBookManager) getOrderBook(ctx context.Context, currencyPair types.CurrencyPair) (*order_book.OrderBook, error) {
	if lo.HasKey(m.orderBooks, currencyPair.String()) {
		return m.orderBooks[currencyPair.String()], nil
	} else if lo.HasKey(m.orderBooks, currencyPair.StringReverse()) {
		return m.orderBooks[currencyPair.StringReverse()], nil
	}

	return nil, echo.ErrNotFound
}

func (m *OrderBookManager) matchOrders(ctx context.Context, orderBook *order_book.OrderBook) error {
	i := 0
	j := 0

	for i < len(orderBook.BuyOrders) && j < len(orderBook.SellOrders) {
		buyOrder := *orderBook.BuyOrders[i]
		sellOrder := *orderBook.SellOrders[j]

		if buyOrder.Status != types.OrderOpen {
			i++
			continue
		} else if sellOrder.Status != types.OrderOpen {
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
				Purpose:     null.StringFrom(fmt.Sprintf("Trading %v for %v", sellOrder.SellCurrency, sellOrder.BuyCurrency)),
			}

			_, err := m.store.Transfer(ctx, transferArgs)
			if err != nil {
				return fmt.Errorf("storefx.Transfer failed: %w", err)
			}

			buyQuantity := sellQuantity * sellOrder.Price

			transferArgs = store_types.TransferArgs{
				FromAddress: buyOrder.SellRequisites.Address,
				ToAddress:   sellOrder.BuyRequisites.Address,
				Amount:      buyQuantity,
				Currency:    sellOrder.BuyCurrency,
				Purpose:     null.StringFrom(fmt.Sprintf("Trading %v for %v", sellOrder.SellCurrency, sellOrder.BuyCurrency)),
			}

			_, err = m.store.Transfer(ctx, transferArgs)
			if err != nil {
				return fmt.Errorf("storefx.Transfer failed: %w", err)
			}

			buyOrder.BuyQuantity = null.FloatFrom(buyOrder.BuyQuantity.ValueOrZero() - sellQuantity)
			sellOrder.SellQuantity = null.FloatFrom(sellOrder.SellQuantity.ValueOrZero() - sellQuantity)

			if buyOrder.BuyQuantity.ValueOrZero() <= 0 {
				buyOrder.Status = types.OrderClosed
				if err = m.store.SaveOrder(ctx, buyOrder); err != nil {
					return fmt.Errorf("storefx.SaveOrder failed: %w", err)
				}
			}
			if sellOrder.SellQuantity.ValueOrZero() <= 0 {
				sellOrder.Status = types.OrderClosed
				if err = m.store.SaveOrder(ctx, sellOrder); err != nil {
					return fmt.Errorf("storefx.SaveOrder failed: %w", err)
				}
			}

			orderBook.BuyOrders[i] = &buyOrder
			orderBook.SellOrders[j] = &sellOrder
		} else {
			break
		}
	}

	return nil
}
