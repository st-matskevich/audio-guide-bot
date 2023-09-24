package controller

import (
	"log"

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

func HandlerPrintf(c *fiber.Ctx, format string, v ...any) {
	args := append([]any{c.Method(), c.Path()}, v...)
	log.Printf("%s %s: "+format, args...)
}

// Handler responses statuses
const (
	RESPONSE_SUCCESS = "success"
	RESPONSE_FAIL    = "fail"
	RESPONSE_ERROR   = "error"
)

type HandlerResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    uint        `json:"code,omitempty"`
}

func HandlerSendSuccess(c *fiber.Ctx, code int, data interface{}) error {
	response := HandlerResponse{
		Status: RESPONSE_SUCCESS,
		Data:   data,
	}
	return c.Status(code).JSON(response)
}

func HandlerSendFailure(c *fiber.Ctx, code int, data interface{}) error {
	response := HandlerResponse{
		Status: RESPONSE_FAIL,
		Data:   data,
	}
	return c.Status(code).JSON(response)
}

func HandlerSendError(c *fiber.Ctx, code int, message string) error {
	response := HandlerResponse{
		Status:  RESPONSE_ERROR,
		Message: message,
	}
	return c.Status(code).JSON(response)
}
