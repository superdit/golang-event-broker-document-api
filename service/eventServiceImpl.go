package service

import (
	"context"
	"event-broker-document-api/config"
	"event-broker-document-api/exception"
	"event-broker-document-api/helper"
	"event-broker-document-api/model"
	"event-broker-document-api/repository"
	"event-broker-document-api/validation"

	"github.com/jinzhu/copier"
	"github.com/lithammer/shortuuid/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type eventServiceImpl struct {
	EventRepository      repository.EventRepository
	ModelRepository      repository.ModelRepository
	BackendRepository    repository.BackendRepository
	SubscriberRepository repository.SubscriberRepository
	Config               config.Config
}

func NewEventService(EventRepository *repository.EventRepository, ModelRepository *repository.ModelRepository,
	BackendRepository *repository.BackendRepository, SubscriberRepository *repository.SubscriberRepository, config config.Config) EventService {
	return &eventServiceImpl{
		EventRepository:      *EventRepository,
		ModelRepository:      *ModelRepository,
		BackendRepository:    *BackendRepository,
		SubscriberRepository: *SubscriberRepository,
		Config:               config,
	}
}

func (service *eventServiceImpl) Insert(request model.EventInsertRequest) (response model.EventInsertResponse) {
	validation.Required(request.BackendPublisherId, "_backend_publisher_id")
	validation.Required(request.ModelId, "_model_id")
	validation.Required(request.Exchange, "exchange")
	// validation.Required(request.TriggerAction, "trigger_action")

	// find backend
	backend := service.BackendRepository.FindById(request.BackendPublisherId)
	if backend.Id == "" {
		helper.Response400("backend id '"+request.BackendPublisherId+"' tidak ditemukan", 0)
	}

	// find model
	model := service.ModelRepository.FindById(request.ModelId)
	if model.Id == "" {
		helper.Response400("model id '"+request.ModelId+"' tidak ditemukan", 0)
	}

	// find existing event
	event := service.EventRepository.FindByExchange(request.Exchange)
	if event.Id != "" {
		helper.Response400("event exchange '"+request.Exchange+"' sudah digunakan", 0)
	}

	backend.TotalExchanges++

	copier.Copy(&event, request)
	event.Id = shortuuid.New()
	event.BackendPublisherName = backend.Name
	event.ModelName = model.Name
	event.ModelUrl = model.Url

	// update total backend
	client, trxOpts := service.EventRepository.SessionTransaction()
	session, err := client.StartSession()
	exception.PanicIfNeeded(err)

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		service.EventRepository.Insert(event, sessionContext)
		service.BackendRepository.Update(backend, sessionContext)
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, trxOpts)
	exception.PanicIfNeeded(err)

	copier.Copy(&response, event)
	helper.MessageOK = "event inserted"
	return response
}

func (service *eventServiceImpl) Update(request model.EventUpdateRequest) (response model.EventUpdateResponse) {
	validation.Required(request.BackendPublisherId, "_backend_publisher_id")
	validation.Required(request.ModelId, "_model_id")
	validation.Required(request.Exchange, "exchange")
	// validation.Required(request.TriggerAction, "trigger_action")

	// find backend
	newBackend := service.BackendRepository.FindById(request.BackendPublisherId)
	if newBackend.Id == "" {
		helper.Response400("new backend id '"+request.BackendPublisherId+"' tidak ditemukan", 0)
	}

	// find model
	model := service.ModelRepository.FindById(request.ModelId)
	if model.Id == "" {
		helper.Response400("model id '"+request.ModelId+"' tidak ditemukan", 0)
	}

	// find existing event exchange name
	event := service.EventRepository.FindExchangeExceptId(request.Exchange, request.Id)
	if event.Id != "" {
		helper.Response400("event exchange '"+request.Exchange+"' sudah digunakan", 0)
	}

	// find existing event
	event = service.EventRepository.FindById(request.Id)
	if event.Id == "" {
		helper.Response400("event id '"+request.Id+"' tidak ditemukan", 0)
	}

	// find existing backend
	oldBackend := service.BackendRepository.FindById(event.BackendPublisherId)
	if oldBackend.Id == "" {
		helper.Response400("old backend id '"+event.BackendPublisherId+"' tidak ditemukan", 0)
	}

	copier.Copy(&event, request)
	event.BackendPublisherName = newBackend.Name
	event.ModelName = model.Name
	event.ModelUrl = model.Url

	// start transactions
	client, trxOpts := service.EventRepository.SessionTransaction()
	session, err := client.StartSession()
	exception.PanicIfNeeded(err)

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		service.EventRepository.Update(event, sessionContext)
		service.SubscriberRepository.UpdateExchangeByEventId(event.Id, event.Exchange)
		if oldBackend.Id != newBackend.Id {
			newBackend.TotalExchanges++
			oldBackend.TotalExchanges--
			service.BackendRepository.Update(newBackend, sessionContext)
			service.BackendRepository.Update(oldBackend, sessionContext)
		}
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, trxOpts)
	exception.PanicIfNeeded(err)
	// end transactions

	copier.Copy(&response, event)
	helper.MessageOK = "event updated"
	return response
}

func (service *eventServiceImpl) Delete(request model.EventDeleteRequest) (response model.EmptyResponse) {
	// find existing event
	event := service.EventRepository.FindById(request.Id)
	if event.Id == "" {
		helper.Response400("event id '"+request.Id+"' tidak ditemukan", 0)
	}

	backend := service.BackendRepository.FindById(event.BackendPublisherId)

	// start transactions
	client, trxOpts := service.EventRepository.SessionTransaction()
	session, err := client.StartSession()
	exception.PanicIfNeeded(err)

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		service.EventRepository.Delete(event.Id, sessionContext)
		if backend.Id != "" {
			backend.TotalExchanges--
			service.BackendRepository.Update(backend, sessionContext)
		}
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, trxOpts)
	exception.PanicIfNeeded(err)
	// end transactions

	helper.MessageOK = "event deleted"
	return response
}

func (service *eventServiceImpl) List(query string) (response model.EventListResponse) {
	models := service.EventRepository.List(query)

	if len(models) > 0 {
		for _, v := range models {
			tmpEntity := model.EventEnityResponse{}

			copier.Copy(&tmpEntity, &v)
			response.Events = append(response.Events, tmpEntity)
			response.Total += 1
		}
	} else {
		response.Events = make([]model.EventEnityResponse, 0)
	}

	helper.MessageOK = "event list retrieved"
	return response
}
