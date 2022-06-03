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

type subscriberServiceImpl struct {
	EventRepository      repository.EventRepository
	SubscriberRepository repository.SubscriberRepository
	BackendRepository    repository.BackendRepository
	ModelRepository      repository.ModelRepository
	Config               config.Config
}

func NewSubscriberService(EventRepository *repository.EventRepository, SubscriberRepository *repository.SubscriberRepository,
	BackendRepository *repository.BackendRepository, ModelRepository *repository.ModelRepository, config config.Config) SubscriberService {
	return &subscriberServiceImpl{
		EventRepository:      *EventRepository,
		SubscriberRepository: *SubscriberRepository,
		BackendRepository:    *BackendRepository,
		ModelRepository:      *ModelRepository,
		Config:               config,
	}
}

func (service *subscriberServiceImpl) Insert(request model.SubscriberInsertRequest) (response model.SubscriberInsertResponse) {
	validation.Required(request.EventId, "_event_id")
	validation.Required(request.BackendSubcriberId, "_backend_subcriber_id")
	validation.Required(request.Queue, "queue")

	// find event
	event := service.EventRepository.FindById(request.EventId)
	if event.Id == "" {
		helper.Response400("event id '"+request.EventId+"' tidak ditemukan", 0)
	}

	// find model by event
	model := service.ModelRepository.FindById(event.ModelId)
	if model.Id == "" {
		helper.Response400("model id '"+event.ModelId+"' tidak ditemukan", 0)
	}

	// find backend
	backend := service.BackendRepository.FindById(request.BackendSubcriberId)
	if backend.Id == "" {
		helper.Response400("backend id '"+request.BackendSubcriberId+"' tidak ditemukan", 0)
	}

	// find existing subscriber
	subscriber := service.SubscriberRepository.FindByQueue(request.Queue)
	if subscriber.Id != "" {
		helper.Response400("subscriber queue '"+request.Queue+"' sudah digunakan", 0)
	}

	event.TotalSubscribers++
	backend.TotalQueues++

	copier.Copy(&subscriber, request)
	subscriber.Id = shortuuid.New()
	subscriber.BackendSubcriberName = backend.Name
	subscriber.Exchange = event.Exchange
	subscriber.ExchangeModelUrl = model.Url

	client, trxOpts := service.EventRepository.SessionTransaction()
	session, err := client.StartSession()
	exception.PanicIfNeeded(err)

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		service.SubscriberRepository.Insert(subscriber, sessionContext)
		service.EventRepository.Update(event, sessionContext)
		service.BackendRepository.Update(backend, sessionContext)
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, trxOpts)
	exception.PanicIfNeeded(err)

	copier.Copy(&response, subscriber)
	helper.MessageOK = "subscriber inserted"
	return response
}

func (service *subscriberServiceImpl) Update(request model.SubscriberUpdateRequest) (response model.SubscriberUpdateResponse) {
	validation.Required(request.EventId, "_event_id")
	validation.Required(request.BackendSubcriberId, "_backend_subcriber_id")
	validation.Required(request.Queue, "queue")

	// find event
	newEvent := service.EventRepository.FindById(request.EventId)
	if newEvent.Id == "" {
		helper.Response400("event id '"+request.EventId+"' tidak ditemukan", 0)
	}

	// find model by event
	newModel := service.ModelRepository.FindById(newEvent.ModelId)
	if newModel.Id == "" {
		helper.Response400("model id '"+newEvent.ModelId+"' tidak ditemukan", 0)
	}

	// find backend
	newBackend := service.BackendRepository.FindById(request.BackendSubcriberId)
	if newBackend.Id == "" {
		helper.Response400("backend id '"+request.BackendSubcriberId+"' tidak ditemukan", 0)
	}

	// find existing subscriber
	subscriber := service.SubscriberRepository.FindQueueExceptId(request.Queue, request.Id)
	if subscriber.Id != "" {
		helper.Response400("subscriber queue '"+request.Queue+"' sudah digunakan", 0)
	}

	// find existing subscriber
	subscriber = service.SubscriberRepository.FindById(request.Id)
	if subscriber.Id == "" {
		helper.Response400("subscriber id '"+subscriber.Id+"' tidak ditemukan", 0)
	}

	// find old event
	oldEvent := service.EventRepository.FindById(subscriber.EventId)
	if oldEvent.Id == "" {
		helper.Response400("old event id '"+subscriber.EventId+"' tidak ditemukan", 0)
	}

	// find old backend
	oldBackend := service.BackendRepository.FindById(subscriber.BackendSubcriberId)
	if oldBackend.Id == "" {
		helper.Response400("old backend id '"+subscriber.BackendSubcriberId+"' tidak ditemukan", 0)
	}

	copier.Copy(&subscriber, request)
	subscriber.BackendSubcriberName = newBackend.Name
	subscriber.Exchange = newEvent.Exchange
	subscriber.ExchangeModelUrl = newModel.Url

	// start transactions
	client, trxOpts := service.SubscriberRepository.SessionTransaction()
	session, err := client.StartSession()
	exception.PanicIfNeeded(err)

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		service.SubscriberRepository.Update(subscriber, sessionContext)
		if oldBackend.Id != newBackend.Id {
			newBackend.TotalExchanges++
			oldBackend.TotalExchanges--
			service.BackendRepository.Update(newBackend, sessionContext)
			service.BackendRepository.Update(oldBackend, sessionContext)
		}
		if newEvent.Id != oldEvent.Id {
			newEvent.TotalSubscribers++
			newEvent.TotalSubscribers--
			service.EventRepository.Update(newEvent, sessionContext)
			service.EventRepository.Update(oldEvent, sessionContext)
		}
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, trxOpts)
	exception.PanicIfNeeded(err)
	// end transactions

	copier.Copy(&response, subscriber)

	helper.MessageOK = "subscriber updated"
	return response
}

func (service *subscriberServiceImpl) Delete(request model.SubscriberDeleteRequest) (response model.EmptyResponse) {
	// find existing subscriber
	subscriber := service.SubscriberRepository.FindById(request.Id)
	if subscriber.Id == "" {
		helper.Response400("subscriber id '"+request.Id+"' tidak ditemukan", 0)
	}

	backend := service.BackendRepository.FindById(subscriber.BackendSubcriberId)
	event := service.EventRepository.FindById(subscriber.EventId)

	// start transactions
	client, trxOpts := service.SubscriberRepository.SessionTransaction()
	session, err := client.StartSession()
	exception.PanicIfNeeded(err)

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		service.SubscriberRepository.Delete(subscriber.Id, sessionContext)
		if backend.Id != "" {
			backend.TotalQueues--
			service.BackendRepository.Update(backend, sessionContext)
		}
		if event.Id != "" {
			event.TotalSubscribers--
			service.EventRepository.Update(event, sessionContext)
		}
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, trxOpts)
	exception.PanicIfNeeded(err)
	// end transactions

	helper.MessageOK = "subscriber deleted"
	return response
}

func (service *subscriberServiceImpl) List(query string) (response model.SubscriberListResponse) {
	models := service.SubscriberRepository.List(query)

	if len(models) > 0 {
		for _, v := range models {
			tmpEntity := model.SubscriberEnityResponse{}

			copier.Copy(&tmpEntity, &v)
			response.Subscribers = append(response.Subscribers, tmpEntity)
			response.Total += 1
		}
	} else {
		response.Subscribers = make([]model.SubscriberEnityResponse, 0)
	}

	helper.MessageOK = "subscriber list retrieved"
	return response
}
