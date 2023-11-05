package controller

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/st-matskevich/audio-guide-bot/api/provider/auth"
	"github.com/st-matskevich/audio-guide-bot/api/provider/blob"
	"github.com/st-matskevich/audio-guide-bot/api/provider/translation"
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
			Path:    "/objects/:code/covers/:index",
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
		HandlerPrintf(c, LOG_ERROR, "Failed to parse authorization token", "error", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse authorization token")
	}

	if !tokenValid {
		HandlerPrintf(c, LOG_WARNING, "Authorization token is invalid")
		return HandlerSendFailure(c, fiber.StatusUnauthorized, "Authorization token is invalid")
	}

	objectCode := c.Params("code")
	language := c.Query("language")
	object, err := controller.getObject(c, objectCode, language)
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to get object", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get object")
	}

	if object == nil {
		HandlerPrintf(c, LOG_WARNING, "Object not found")
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
		HandlerPrintf(c, LOG_ERROR, "Failed to parse authorization token", "error", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse authorization token")
	}

	if !tokenValid {
		HandlerPrintf(c, LOG_WARNING, "Authorization token is invalid")
		return HandlerSendFailure(c, fiber.StatusUnauthorized, "Authorization token is invalid")
	}

	objectCode := c.Params("code")
	language := c.Query("language")
	object, err := controller.getObject(c, objectCode, language)
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to get object", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get object")
	}

	if object == nil {
		HandlerPrintf(c, LOG_WARNING, "Object not found")
		return HandlerSendFailure(c, fiber.StatusNotFound, "Object not found")
	}

	coverIndex, err := strconv.Atoi(c.Params("index"))
	if err != nil {
		HandlerPrintf(c, LOG_WARNING, "Failed to parse cover index", "error", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse cover index")
	}

	coverPath := ""
	for _, cover := range object.Covers {
		if cover.Index == coverIndex {
			coverPath = cover.Path
			break
		}
	}

	if coverPath == "" {
		HandlerPrintf(c, LOG_WARNING, "Cover not found")
		return HandlerSendFailure(c, fiber.StatusNotFound, "Cover not found")
	}

	reader, err := controller.BlobProvider.ReadBlob(coverPath, blob.ReadBlobOptions{})
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Blob read failed", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Blob read failed")
	}

	c.Type(filepath.Ext(coverPath))
	c.Status(fiber.StatusOK)

	return c.SendStream(reader)
}

func (controller *ObjectsController) HandleGetObjectAudio(c *fiber.Ctx) error {
	// Resources are loaded by HTML components thus it's not possible to pass an access token as a header
	// URL parameter have to be used instead
	authHeader := c.Query("access-token")
	_, tokenValid, err := controller.TokenProvider.Verify(authHeader)

	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to parse authorization token", "error", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse authorization token")
	}

	if !tokenValid {
		HandlerPrintf(c, LOG_WARNING, "Authorization token is invalid")
		return HandlerSendFailure(c, fiber.StatusUnauthorized, "Authorization token is invalid")
	}

	objectCode := c.Params("code")
	language := c.Query("language")
	object, err := controller.getObject(c, objectCode, language)
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to get object", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get object")
	}

	if object == nil {
		HandlerPrintf(c, LOG_WARNING, "Object not found")
		return HandlerSendFailure(c, fiber.StatusNotFound, "Object not found")
	}

	rangesHeader := c.Get(fiber.HeaderRange)
	if rangesHeader != "" {
		blobStat, err := controller.BlobProvider.StatBlob(object.AudioPath)
		if err != nil {
			HandlerPrintf(c, LOG_ERROR, "Blob stat failed", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Blob stat failed")
		}

		units, ranges, err := controller.parseRange(rangesHeader, blobStat.Size)
		if err != nil {
			HandlerPrintf(c, LOG_WARNING, "Failed to parse range header", "error", err)
			return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse range header")
		}

		if units != "bytes" {
			HandlerPrintf(c, LOG_WARNING, "Incorrect range units", "error", units)
			return HandlerSendFailure(c, fiber.StatusRequestedRangeNotSatisfiable, "Incorrect range units")
		}

		if len(ranges) < 1 {
			HandlerPrintf(c, LOG_WARNING, "No ranges provided")
			return HandlerSendFailure(c, fiber.StatusRequestedRangeNotSatisfiable, "No ranges provided")
		}

		readOptions := blob.ReadBlobOptions{}
		readOptions.Range = &ranges[0]

		reader, err := controller.BlobProvider.ReadBlob(object.AudioPath, readOptions)
		if err != nil {
			HandlerPrintf(c, LOG_ERROR, "Blob read failed", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Blob read failed")
		}

		c.Type(filepath.Ext(object.AudioPath))
		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", readOptions.Range.Start, readOptions.Range.End, blobStat.Size))
		c.Set("Accept-Ranges", "bytes")
		c.Status(fiber.StatusPartialContent)

		return c.SendStream(reader)
	} else {
		reader, err := controller.BlobProvider.ReadBlob(object.AudioPath, blob.ReadBlobOptions{})
		if err != nil {
			HandlerPrintf(c, LOG_ERROR, "Blob read failed", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Blob read failed")
		}

		c.Type(filepath.Ext(object.AudioPath))
		c.Set("Accept-Ranges", "bytes")
		c.Status(fiber.StatusOK)

		return c.SendStream(reader)
	}
}

func (controller *ObjectsController) getObject(c *fiber.Ctx, code string, language string) (*repository.Object, error) {
	object, err := controller.ObjectRepository.GetObject(code, language)
	if err != nil {
		return nil, err
	}

	if object == nil {
		// fallback to default language
		fallback := translation.DEFAULT_LANGUAGE.String()
		HandlerPrintf(c, LOG_WARNING, "Object i18n for requested language is not found, loading i18n for fallback language", "language", language, "fallback", fallback)
		object, err = controller.ObjectRepository.GetObject(code, fallback)
		if err != nil {
			return nil, err
		}
	}

	return object, nil
}

func (controller *ObjectsController) parseRange(header string, size int64) (string, []blob.BlobRange, error) {
	if header == "" || !strings.Contains(header, "=") {
		return "", nil, errors.New("malformed range header string")
	}

	data := strings.Split(header, "=")
	const expectedDataParts = 2
	if len(data) != expectedDataParts {
		return "", nil, errors.New("malformed range header string")
	}

	units := data[0]
	ranges := []blob.BlobRange{}
	arr := strings.Split(data[1], ",")
	for i := 0; i < len(arr); i++ {
		item := strings.Split(arr[i], "-")
		if len(item) == 1 {
			return "", nil, errors.New("malformed range header string")
		}
		start, startErr := strconv.ParseInt(item[0], 10, 64)
		end, endErr := strconv.ParseInt(item[1], 10, 64)
		if startErr != nil { // -nnn
			start = size - end
			end = size - 1
		} else if endErr != nil { // nnn-
			end = size - 1
		}
		if end > size-1 { // limit last-byte-pos to current length
			end = size - 1
		}
		if start > end || start < 0 {
			continue
		}
		ranges = append(ranges, blob.BlobRange{
			Start: start,
			End:   end,
		})
	}

	return units, ranges, nil
}
