package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"vitalik_backend/internal/app"
)

type Server struct {
	logger *zap.SugaredLogger
	echo   *echo.Echo
	app    *app.Application
}

func NewServer(app *app.Application) *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	s := &Server{
		logger: app.Logger,
		app:    app,
		echo:   e,
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
	s.echo.GET("/healthCheck", s.app.HealthCheck())
}
