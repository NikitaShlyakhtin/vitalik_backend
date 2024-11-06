package app

import (
	"go.uber.org/zap"
)

// Application holds the application state and dependencies
type Application struct {
	Logger *zap.SugaredLogger
}

// NewApplication initializes a new Application instance
func NewApplication(
	logger *zap.Logger,
) *Application {
	return &Application{
		Logger: logger.Sugar(),
	}
}
