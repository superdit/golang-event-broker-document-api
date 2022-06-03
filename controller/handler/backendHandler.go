package handler

import (
	"event-broker-document-api/exception"
	"event-broker-document-api/helper"
	"event-broker-document-api/model"
	"event-broker-document-api/service"

	"github.com/gofiber/fiber/v2"
)

type BackendHandler struct {
	BackendService service.BackendService
}

func NewBackendHandler(
	BackendService *service.BackendService) BackendHandler {

	return BackendHandler{
		BackendService: *BackendService,
	}
}

func (handler *BackendHandler) Insert(c *fiber.Ctx) error {
	var request model.BackendInsertRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)

	response := handler.BackendService.Insert(request)
	return helper.ResponseOK(c, response)
}

func (handler *BackendHandler) Update(c *fiber.Ctx) error {
	var request model.BackendUpdateRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)
	request.Id = c.Params("id")

	response := handler.BackendService.Update(request)
	return helper.ResponseOK(c, response)
}

func (handler *BackendHandler) Delete(c *fiber.Ctx) error {
	var request model.BackendDeleteRequest

	request.Id = c.Params("id")

	response := handler.BackendService.Delete(request)
	return helper.ResponseOK(c, response)
}

func (handler *BackendHandler) List(c *fiber.Ctx) error {
	query := c.Query("q")
	response := handler.BackendService.List(query)
	return helper.ResponseOK(c, response)
}
