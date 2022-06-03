package main

import (
	"crypto/tls"
	"event-broker-document-api/config"
	"event-broker-document-api/controller"
	"event-broker-document-api/exception"
	"event-broker-document-api/repository"
	"event-broker-document-api/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	configuration := config.New()
	database := config.NewMongoDatabase(configuration)

	app := fiber.New(config.NewFiberConfig())
	// Default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(recover.New())

	// Setup Repository
	modelRepository := repository.NewModelRepository(database)
	backendRepository := repository.NewBackendRepository(database)
	eventRespository := repository.NewEventRepository(database)
	subscriberRepository := repository.NewSubscriberRepository(database)

	// Setup Service
	modelService := service.NewModelService(&modelRepository, &eventRespository, &subscriberRepository, configuration)
	backendService := service.NewBackendService(&backendRepository, &eventRespository, &subscriberRepository, configuration)
	eventService := service.NewEventService(&eventRespository, &modelRepository, &backendRepository, &subscriberRepository, configuration)
	subscriberService := service.NewSubscriberService(&eventRespository, &subscriberRepository, &backendRepository, &modelRepository, configuration)

	// Setup Controller
	appController := controller.NewAppController(&modelService, &backendService, &eventService, &subscriberService)

	// Setup Routing
	appController.Route(app)

	app.Static("/", "./public")
	app.Static("/backend*", "./public")
	app.Static("/model*", "./public")
	app.Static("/event*", "./public")
	app.Static("/subscriber*", "./public")

	// Start App
	if configuration.Get("SSL") == "1" {
		cer, err := tls.LoadX509KeyPair(configuration.Get("SSL_CERT"), configuration.Get("SSL_CERT_KEY"))
		exception.PanicIfNeeded(err)

		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		// Create custom listener
		ln, err := tls.Listen("tcp", ":"+configuration.Get("SSL_PORT"), config)
		exception.PanicIfNeeded(err)

		log.Fatal(app.Listener(ln))
	}

	err := app.Listen(":3000")
	exception.PanicIfNeeded(err)
}
