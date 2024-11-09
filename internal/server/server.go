package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"vitalik_backend/internal/dependencies"
)

type Server struct {
	logger  *zap.SugaredLogger
	echo    *echo.Echo
	handler dependencies.IHandler
}

func NewServer(
	logger *zap.Logger,
	handler dependencies.IHandler,
) *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	s := &Server{
		logger:  logger.Sugar(),
		echo:    e,
		handler: handler,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) Start(address string) error {
	s.logger.Infof("starting on address: %v", address)

	return s.echo.Start(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down gracefully, press Ctrl+C again to force")

	return s.echo.Shutdown(ctx)
}

func (s *Server) setupMiddleware() {
	s.echo.Use(requestLogger(s.logger))
	s.echo.Use(middleware.Recover())
}

func (s *Server) setupRoutes() {
	s.echo.GET("/healthCheck", s.handler.HealthCheck())

	s.echo.POST("/wallets/create", s.handler.CreateWallet())
	s.echo.GET("/wallets", s.handler.ListWallets())
	s.echo.PUT("/wallets/deposit", s.handler.Deposit())

	s.echo.POST("/transfer", s.handler.Transfer())

	s.echo.GET("/transactions", s.handler.ListTransactions())

	s.echo.POST("/orders/create", s.handler.CreateOrder())
	s.echo.DELETE("/orders", s.handler.CancelOrder())
	s.echo.GET("/orders", s.handler.ListOrders())
	s.echo.POST("/orders/match", s.handler.MatchOrders())
}
