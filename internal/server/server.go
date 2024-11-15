package server

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"vitalik_backend/internal/dependencies"
)

type Server struct {
	logger  *zap.SugaredLogger
	echo    *echo.Echo
	handler dependencies.IHandler
	auth    dependencies.IAuthService
}

func NewServer(
	logger *zap.Logger,
	handler dependencies.IHandler,
	auth dependencies.IAuthService,
) (*Server, error) {
	if logger == nil ||
		handler == nil ||
		auth == nil {
		return nil, errors.New("failed to initialize server")
	}

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	s := &Server{
		logger:  logger.Sugar(),
		echo:    e,
		handler: handler,
		auth:    auth,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
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

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
}

func (s *Server) setupRoutes() {
	s.echo.GET("/healthCheck", s.handler.HealthCheck())
	s.echo.POST("/register", s.handler.Register())
	s.echo.POST("/login", s.handler.Login())

	authGroup := s.echo.Group("/auth")
	authGroup.Use(s.authMiddleware)

	authGroup.POST("/wallets/create", s.handler.CreateWallet())
	authGroup.POST("/wallets", s.handler.ListWallets())
	authGroup.PUT("/wallets/deposit", s.handler.Deposit())

	authGroup.POST("/transfer", s.handler.Transfer())
	authGroup.POST("/transactions", s.handler.ListTransactions())
	authGroup.POST("/orders/create", s.handler.CreateOrder())
	authGroup.DELETE("/orders", s.handler.CancelOrder())
	authGroup.POST("/orders", s.handler.ListOrders())
	authGroup.POST("/orders/match", s.handler.MatchOrders())

	authGroup.GET("/currencies", s.handler.ListAvailableCurrencies())
}
