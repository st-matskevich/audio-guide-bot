package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/st-matskevich/audio-guide-bot/api/auth"
	"github.com/st-matskevich/audio-guide-bot/api/blob"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

type ObjectsController struct {
	TokenProvider    auth.TokenProvider
	BlobProvider     blob.BlobProvider
	ObjectRepository repository.ObjectRepository
}

func (controller *ObjectsController) GetRoutes() []Route {
	return []Route{
		{
			Method:  "GET",
			Path:    "/objects/:code",
			Handler: controller.HandleGetObject,
		},
		{
			Method:  "GET",
			Path:    "/objects/:code/cover",
			Handler: controller.HandleGetObjectCover,
		},
		{
			Method:  "GET",
			Path:    "/objects/:code/audio",
			Handler: controller.HandleGetObjectAudio,
		},
	}
}

func (controller *ObjectsController) HandleGetObject(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	_, tokenValid, err := controller.TokenProvider.Verify(authHeader)

	if err != nil {
		HandlerPrintf(c, "Failed to parse authorization token - %v", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse authorization token")
	}

	if !tokenValid {
		HandlerPrintf(c, "Authorization token is invalid")
		return HandlerSendFailure(c, fiber.StatusUnauthorized, "Authorization token is invalid")
	}

	objectCode := c.Params("code")
	object, err := controller.ObjectRepository.GetObject(objectCode)
	if err != nil {
		HandlerPrintf(c, "Failed to get object - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get object")
	}

	if object == nil {
		HandlerPrintf(c, "Object not found")
		return HandlerSendFailure(c, fiber.StatusNotFound, "Object not found")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, object)
}

func (controller *ObjectsController) HandleGetObjectCover(c *fiber.Ctx) error {
	// Resources are loaded by HTML components thus it's not possible to pass an access token as a header
	// URL parameter have to be used instead
	authHeader := c.Query("access-token")
	_, tokenValid, err := controller.TokenProvider.Verify(authHeader)

	if err != nil {
		HandlerPrintf(c, "Failed to parse authorization token - %v", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse authorization token")
	}

	if !tokenValid {
		HandlerPrintf(c, "Authorization token is invalid")
		return HandlerSendFailure(c, fiber.StatusUnauthorized, "Authorization token is invalid")
	}

	objectCode := c.Params("code")
	object, err := controller.ObjectRepository.GetObject(objectCode)
	if err != nil {
		HandlerPrintf(c, "Failed to get object - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get object")
	}

	if object == nil {
		HandlerPrintf(c, "Object not found")
		return HandlerSendFailure(c, fiber.StatusNotFound, "Object not found")
	}

	err = controller.BlobProvider.ReadBlob(object.CoverPath, c.Response().BodyWriter())
	if err != nil {
		HandlerPrintf(c, "Blob read failed - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Blob read failed")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (controller *ObjectsController) HandleGetObjectAudio(c *fiber.Ctx) error {
	// Resources are loaded by HTML components thus it's not possible to pass an access token as a header
	// URL parameter have to be used instead
	authHeader := c.Query("access-token")
	_, tokenValid, err := controller.TokenProvider.Verify(authHeader)

	if err != nil {
		HandlerPrintf(c, "Failed to parse authorization token - %v", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse authorization token")
	}

	if !tokenValid {
		HandlerPrintf(c, "Authorization token is invalid")
		return HandlerSendFailure(c, fiber.StatusUnauthorized, "Authorization token is invalid")
	}

	objectCode := c.Params("code")
	object, err := controller.ObjectRepository.GetObject(objectCode)
	if err != nil {
		HandlerPrintf(c, "Failed to get object - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get object")
	}

	if object == nil {
		HandlerPrintf(c, "Object not found")
		return HandlerSendFailure(c, fiber.StatusNotFound, "Object not found")
	}

	err = controller.BlobProvider.ReadBlob(object.AudioPath, c.Response().BodyWriter())
	if err != nil {
		HandlerPrintf(c, "Blob read failed - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Blob read failed")
	}

	// TODO: Accept-Ranges header allows seeking in the player,
	// but proper Range header handling still needs to be implemented
	c.Set("Accept-Ranges", "bytes")
	return c.SendStatus(fiber.StatusOK)
}
