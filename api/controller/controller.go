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
