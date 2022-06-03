package handler

import (
	"event-broker-document-api/exception"
	"event-broker-document-api/helper"
	"event-broker-document-api/model"
	"event-broker-document-api/service"

	"github.com/gofiber/fiber/v2"
)

type SubscriberHandler struct {
	SubscriberService service.SubscriberService
}

func NewSubscriberHandler(
	SubscriberService *service.SubscriberService) SubscriberHandler {

	return SubscriberHandler{
		SubscriberService: *SubscriberService,
	}
}

func (handler *SubscriberHandler) Insert(c *fiber.Ctx) error {
	var request model.SubscriberInsertRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)

	response := handler.SubscriberService.Insert(request)
	return helper.ResponseOK(c, response)
}

func (handler *SubscriberHandler) Update(c *fiber.Ctx) error {
	var request model.SubscriberUpdateRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)
	request.Id = c.Params("id")

	response := handler.SubscriberService.Update(request)
	return helper.ResponseOK(c, response)
}

func (handler *SubscriberHandler) Delete(c *fiber.Ctx) error {
	var request model.SubscriberDeleteRequest

	request.Id = c.Params("id")

	response := handler.SubscriberService.Delete(request)
	return helper.ResponseOK(c, response)
}

func (handler *SubscriberHandler) List(c *fiber.Ctx) error {
	query := c.Query("q")
	response := handler.SubscriberService.List(query)
	return helper.ResponseOK(c, response)
}
