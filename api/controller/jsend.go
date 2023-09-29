package controller

import "github.com/gofiber/fiber/v2"

// JSend response statuses
const (
	RESPONSE_SUCCESS = "success"
	RESPONSE_FAIL    = "fail"
	RESPONSE_ERROR   = "error"
)

// JSend response data
type HandlerResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    uint        `json:"code,omitempty"`
}

// JSend response handlers
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
