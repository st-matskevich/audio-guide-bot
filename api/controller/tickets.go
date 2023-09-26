package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

type TicketsController struct {
	JWTSecret        []byte
	TicketRepository repository.TicketRepository
}

func (controller *TicketsController) GetRoutes() []Route {
	return []Route{
		{
			Method:  "POST",
			Path:    "/tickets/token",
			Handler: controller.HandleExchangeTicketForToken,
		},
	}
}

func (controller *TicketsController) HandleExchangeTicketForToken(c *fiber.Ctx) error {
	input := struct {
		Ticket uuid.UUID `json:"ticket"`
	}{}

	if err := c.BodyParser(&input); err != nil {
		HandlerPrintf(c, "Failed to parse input - %v", err)
		return HandlerSendError(c, fiber.StatusBadRequest, "Failed to parse input")
	}

	active, err := controller.TicketRepository.ActivateTicket(input.Ticket.String())
	if err != nil {
		HandlerPrintf(c, "Failed to activate ticket - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to activate ticket")
	}

	if !active {
		HandlerPrintf(c, "Requested ticket already activated")
		return HandlerSendFailure(c, fiber.StatusForbidden, "Requested ticket already activated")
	}

	currentTime := time.Now()
	expires := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, 0, 0, 0, 0, time.UTC)
	claims := jwt.MapClaims{
		"exp": expires.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(controller.JWTSecret)
	if err != nil {
		HandlerPrintf(c, "Failed to sign token - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to sign token")
	}

	result := struct {
		Token string `json:"token"`
	}{}

	result.Token = tokenString
	return HandlerSendSuccess(c, fiber.StatusCreated, result)
}
