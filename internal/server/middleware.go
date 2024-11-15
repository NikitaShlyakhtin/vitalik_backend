package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func requestLogger(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogUserAgent: true,
		LogRemoteIP:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Status >= 200 && v.Status < 300 {
				logger.Infow("REQUEST",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"user_agent", v.UserAgent,
					"remote_ip", v.RemoteIP,
				)
			} else {
				logger.Infow("REQUEST_ERROR",
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

func (s *Server) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
		}

		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
		}
		tokenString := tokenParts[1]

		userID, err := s.auth.Authenticate(c.Request().Context(), tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set("user_id", userID)

		return next(c)
	}
}
