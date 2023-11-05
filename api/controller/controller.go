package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetRoutes() []Route
}

type Route struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

const (
	LOG_INFO    = 0
	LOG_WARNING = 1
	LOG_ERROR   = 2
)

func HandlerPrintf(c *fiber.Ctx, severity int, message string, v ...any) {
	logger := slog.Info
	switch severity {
	case LOG_WARNING:
		logger = slog.Warn
	case LOG_ERROR:
		logger = slog.Error
	}

	args := append([]any{GetRequestLogGroup(c)}, v...)
	logger(message, args...)
}
