package controller

import (
	"event-broker-document-api/controller/handler"
	"event-broker-document-api/service"

	"github.com/gofiber/fiber/v2"
)

type AppController struct {
	ModelService      service.ModelService
	BackendService    service.BackendService
	EventService      service.EventService
	SubscriberService service.SubscriberService
}

func NewAppController(ModelService *service.ModelService, BackendService *service.BackendService,
	EventService *service.EventService, SubscriberService *service.SubscriberService) AppController {

	return AppController{
		ModelService:      *ModelService,
		BackendService:    *BackendService,
		EventService:      *EventService,
		SubscriberService: *SubscriberService,
	}
}

func (controller *AppController) Route(app *fiber.App) {

	authHandler := handler.NewModelHandler(&controller.ModelService)
	app.Post("/api/model", authHandler.Insert)
	app.Put("/api/model/:id", authHandler.Update)
	app.Delete("/api/model/:id", authHandler.Delete)
	app.Get("/api/models", authHandler.List)

	backendHandler := handler.NewBackendHandler(&controller.BackendService)
	app.Post("/api/backend", backendHandler.Insert)
	app.Put("/api/backend/:id", backendHandler.Update)
	app.Delete("/api/backend/:id", backendHandler.Delete)
	app.Get("/api/backends", backendHandler.List)

	eventHandler := handler.NewEventHandler(&controller.EventService)
	app.Post("/api/event", eventHandler.Insert)
	app.Put("/api/event/:id", eventHandler.Update)
	app.Delete("/api/event/:id", eventHandler.Delete)
	app.Get("/api/events", eventHandler.List)

	subscriberHandler := handler.NewSubscriberHandler(&controller.SubscriberService)
	app.Post("/api/subscriber", subscriberHandler.Insert)
	app.Put("/api/subscriber/:id", subscriberHandler.Update)
	app.Delete("/api/subscriber/:id", subscriberHandler.Delete)
	app.Get("/api/subscribers", subscriberHandler.List)
}
