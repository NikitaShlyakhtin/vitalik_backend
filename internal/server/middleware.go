package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func requestLogger(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogUserAgent: true,
		LogRemoteIP:  true,
		LogError:     true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Infow("REQUEST",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"user_agent", v.UserAgent,
					"remote_ip", v.RemoteIP,
				)
			} else {
				logger.Errorw("REQUEST_ERROR",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"user_agent", v.UserAgent,
					"remote_ip", v.RemoteIP,
					"error", v.Error,
				)
			}

			return nil
		},
	})
}
