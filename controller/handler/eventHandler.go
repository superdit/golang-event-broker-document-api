package handler

import (
	"event-broker-document-api/exception"
	"event-broker-document-api/helper"
	"event-broker-document-api/model"
	"event-broker-document-api/service"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	EventService service.EventService
}

func NewEventHandler(
	EventService *service.EventService) EventHandler {

	return EventHandler{
		EventService: *EventService,
	}
}

func (handler *EventHandler) Insert(c *fiber.Ctx) error {
	var request model.EventInsertRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)

	response := handler.EventService.Insert(request)
	return helper.ResponseOK(c, response)
}

func (handler *EventHandler) Update(c *fiber.Ctx) error {
	var request model.EventUpdateRequest
	err := c.BodyParser(&request)

	exception.PanicIfBadRequest(err)
	request.Id = c.Params("id")

	response := handler.EventService.Update(request)
	return helper.ResponseOK(c, response)
}

func (handler *EventHandler) Delete(c *fiber.Ctx) error {
	var request model.EventDeleteRequest

	request.Id = c.Params("id")

	response := handler.EventService.Delete(request)
	return helper.ResponseOK(c, response)
}

func (handler *EventHandler) List(c *fiber.Ctx) error {
	query := c.Query("q")
	response := handler.EventService.List(query)
	return helper.ResponseOK(c, response)
}
