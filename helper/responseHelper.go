package helper

import (
	"event-broker-document-api/exception"
	"event-broker-document-api/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var MessageOK string

func ResponseExit(httpCode int, message string) {
	exception.PanicIfNeeded(exception.ValidationError{Status: httpCode, Message: message})
}

func ResponseError(httpCode int, message string, errorCode int) {
	exception.PanicIfNeeded(exception.ValidationError{Status: httpCode, Message: message, ErrorCode: errorCode})
}

func Response400(message string, errorCode int) {
	exception.PanicIfNeeded(exception.ValidationError{Status: 400, Message: message, ErrorCode: errorCode})
}

func Response401(message string, errorCode int) {
	exception.PanicIfNeeded(exception.ValidationError{Status: 401, Message: message, ErrorCode: errorCode})
}

func ResponseOK(c *fiber.Ctx, response interface{}) error {

	return c.Status(http.StatusOK).JSON(model.WebResponse{
		Data:    response,
		Message: MessageOK,
	})
}
