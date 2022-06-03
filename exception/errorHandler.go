package exception

import (
	"event-broker-document-api/model"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	status, _ := err.(ValidationError)

	switch status.Status {
	case 400:
		status.Message = "BAD_REQUEST"
	case 401:
		status.Message = "UNAUTHORIZED"
	default:
		status.Status = 500
		status.Message = "INTERNAL_SERVER_ERROR"
	}

	if e, ok := err.(*fiber.Error); ok {
		// handle 405 method not allowed
		if e.Code == 405 {
			status.Status, status.ErrorCode = e.Code, e.Code
		}
	}

	return ctx.Status(status.Status).JSON(model.WebResponseError{
		Error: model.ErrorResponse{
			Message: err.Error(),
			Code:    status.ErrorCode,
		},
	})
}
