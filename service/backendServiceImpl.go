package service

import (
	"event-broker-document-api/config"
	"event-broker-document-api/helper"
	"event-broker-document-api/model"
	"event-broker-document-api/repository"
	"event-broker-document-api/validation"

	"github.com/jinzhu/copier"
	"github.com/lithammer/shortuuid/v3"
)

type backendServiceImpl struct {
	BackendRepository    repository.BackendRepository
	EventRepository      repository.EventRepository
	SubscriberRepository repository.SubscriberRepository
	Config               config.Config
}

func NewBackendService(BackendRepository *repository.BackendRepository, EventRepository *repository.EventRepository,
	SubscriberRepository *repository.SubscriberRepository, config config.Config) BackendService {
	return &backendServiceImpl{
		BackendRepository:    *BackendRepository,
		EventRepository:      *EventRepository,
		SubscriberRepository: *SubscriberRepository,
		Config:               config,
	}
}

func (service *backendServiceImpl) Insert(request model.BackendInsertRequest) (response model.BackendInsertResponse) {
	validation.Required(request.Name, "name")

	backend := service.BackendRepository.FindByName(request.Name)
	if backend.Id != "" {
		helper.ResponseError(400, "nama backend sudah digunakan", 0)
	}

	copier.Copy(&backend, request)
	backend.Id = shortuuid.New()

	service.BackendRepository.Insert(backend)
	copier.Copy(&response, backend)

	helper.MessageOK = "backend inserted"

	return response
}

func (service *backendServiceImpl) Update(request model.BackendUpdateRequest) (response model.BackendUpdateResponse) {
	validation.Required(request.Name, "name")

	backend := service.BackendRepository.FindNameExceptId(request.Name, request.Id)
	if backend.Id != "" {
		helper.Response400("nama backend sudah digunakan", 0)
	}

	backend = service.BackendRepository.FindById(request.Id)
	if backend.Id == "" {
		helper.Response400("id backend tidak ditemukan", 0)
	}

	copier.Copy(&backend, request)
	service.BackendRepository.Update(backend)

	// semua backend name harus berubah di event & subscribers
	service.SubscriberRepository.UpdateBackendByBackendId(backend.Id, backend.Name)
	service.EventRepository.UpdateBackendByBackendId(backend.Id, backend.Name)

	copier.Copy(&response, backend)
	helper.MessageOK = "backend updated"

	return response
}

func (service *backendServiceImpl) Delete(request model.BackendDeleteRequest) (response model.EmptyResponse) {
	model := service.BackendRepository.FindById(request.Id)
	if model.Id == "" {
		helper.Response400("id backend tidak ditemukan", 0)
	}

	service.BackendRepository.Delete(request.Id)
	helper.MessageOK = "backend deleted"

	return response
}

func (service *backendServiceImpl) List(query string) (response model.BackendListResponse) {
	models := service.BackendRepository.List(query)

	if len(models) > 0 {
		for _, v := range models {
			tmpEntity := model.BackendEnityResponse{}

			copier.Copy(&tmpEntity, &v)
			response.Backends = append(response.Backends, tmpEntity)
			response.Total += 1
		}
	} else {
		response.Backends = make([]model.BackendEnityResponse, 0)
	}

	helper.MessageOK = "backend list retrieved"
	return response
}
