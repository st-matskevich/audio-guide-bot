package controller

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Replace logs keys to allow GCP to parse logs correctly
// https://cloud.google.com/logging/docs/structured-logging
func CreateGCPLoggerAdapter() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case "level":
			a.Key = "severity"
		case "msg":
			a.Key = "message"
		}
		return a
	}
}

const CTX_REQUEST_ID = "requestID"

func CreateLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Context().SetUserValue(CTX_REQUEST_ID, requestID)
		c.Set("X-Request-ID", requestID)

		start := time.Now()

		err := c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := c.Response().StatusCode()
		userAgent := string(c.Context().UserAgent())

		logger := slog.Info
		message := "Responeded to request with success code"
		if status >= http.StatusInternalServerError {
			logger = slog.Error
			message = "Responeded to request with error code"
		} else if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			logger = slog.Warn
			message = "Responeded to request with failure code"
		}

		logger(message,
			"userAgent", userAgent,
			"status", status,
			"latency", latency,
			GetRequestLogGroup(c),
		)

		return err
	}
}

func GetRequestLogGroup(c *fiber.Ctx) slog.Attr {
	id := GetRequestID(c)
	return slog.Group("httpRequest",
		"id", id,
		"method", c.Method(),
		"path", c.Path(),
	)
}

func GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Context().UserValue(CTX_REQUEST_ID).(string)
	if !ok {
		return ""
	}

	return requestID
}
