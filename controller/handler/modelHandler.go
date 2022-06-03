package handler

import (
	"event-broker-document-api/exception"
	"event-broker-document-api/helper"
	"event-broker-document-api/model"
	"event-broker-document-api/service"

	"github.com/gofiber/fiber/v2"
)

type ModelHandler struct {
	ModelService service.ModelService
}

func NewModelHandler(
	ModelService *service.ModelService) ModelHandler {

	return ModelHandler{
		ModelService: *ModelService,
	}
}

func (handler *ModelHandler) Insert(c *fiber.Ctx) error {
	var request model.ModelInsertRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)

	response := handler.ModelService.Insert(request)
	return helper.ResponseOK(c, response)
}

func (handler *ModelHandler) Update(c *fiber.Ctx) error {
	var request model.ModelUpdateRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)
	request.Id = c.Params("id")

	response := handler.ModelService.Update(request)
	return helper.ResponseOK(c, response)
}

func (handler *ModelHandler) Delete(c *fiber.Ctx) error {
	var request model.ModelDeleteRequest

	request.Id = c.Params("id")

	response := handler.ModelService.Delete(request)
	return helper.ResponseOK(c, response)
}

func (handler *ModelHandler) List(c *fiber.Ctx) error {
	query := c.Query("q")
	response := handler.ModelService.List(query)
	return helper.ResponseOK(c, response)
}
