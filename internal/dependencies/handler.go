package dependencies

import "github.com/labstack/echo/v4"

// IHandler defines the methods for HTTP handlers
type IHandler interface {
	HealthCheck() echo.HandlerFunc

	CreateWallet() echo.HandlerFunc
	ListWallets() echo.HandlerFunc
	Deposit() echo.HandlerFunc
	ListTransactions() echo.HandlerFunc
	Transfer() echo.HandlerFunc
	CreateOrder() echo.HandlerFunc
	CancelOrder() echo.HandlerFunc
	ListOrders() echo.HandlerFunc
	MatchOrders() echo.HandlerFunc
	ListAvailableCurrencies() echo.HandlerFunc
}
