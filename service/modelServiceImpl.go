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

type modelServiceImpl struct {
	ModelRepository      repository.ModelRepository
	EventRepository      repository.EventRepository
	SubscriberRepository repository.SubscriberRepository
	Config               config.Config
}

func NewModelService(ModelRepository *repository.ModelRepository, EventRepository *repository.EventRepository,
	SubscriberRepository *repository.SubscriberRepository, config config.Config) ModelService {
	return &modelServiceImpl{
		ModelRepository:      *ModelRepository,
		EventRepository:      *EventRepository,
		SubscriberRepository: *SubscriberRepository,
		Config:               config,
	}
}

func (service *modelServiceImpl) Insert(request model.ModelInsertRequest) (response model.ModelInsertResponse) {

	validation.Required(request.Name, "name")
	validation.Required(request.Url, "url")

	model := service.ModelRepository.FindByName(request.Name)
	if model.Id != "" {
		helper.Response400("nama model sudah digunakan", 0)
	}

	copier.Copy(&model, request)
	model.Id = shortuuid.New()

	service.ModelRepository.Insert(model)
	copier.Copy(&response, model)

	helper.MessageOK = "model inserted"

	return response
}

func (service *modelServiceImpl) Update(request model.ModelUpdateRequest) (response model.ModelUpdateResponse) {

	validation.Required(request.Name, "name")
	validation.Required(request.Url, "url")

	model := service.ModelRepository.FindById(request.Id)
	if model.Id == "" {
		helper.Response400("id model tidak ditemukan", 0)
	}

	model = service.ModelRepository.FindNameExceptId(request.Name, request.Id)
	if model.Id != "" {
		helper.Response400("nama model sudah digunakan", 0)
	}

	copier.Copy(&model, request)
	service.ModelRepository.Update(model)

	// url model di event dan subscriber harus diubah
	service.EventRepository.UpdateModelByModelId(model.Id, model.Name, model.Url)

	copier.Copy(&response, model)
	helper.MessageOK = "model updated"

	return response
}

func (service *modelServiceImpl) Delete(request model.ModelDeleteRequest) (response model.EmptyResponse) {

	model := service.ModelRepository.FindById(request.Id)
	if model.Id == "" {
		helper.Response400("id model tidak ditemukan", 0)
	}

	service.ModelRepository.Delete(request.Id)
	helper.MessageOK = "model deleted"
	return response
}

func (service *modelServiceImpl) List(query string) (response model.ModelListResponse) {

	models := service.ModelRepository.List(query)

	if len(models) > 0 {
		for _, v := range models {
			tmpEntity := model.ModelEnityResponse{}

			copier.Copy(&tmpEntity, &v)
			response.Models = append(response.Models, tmpEntity)
			response.Total += 1
		}
	} else {
		response.Models = make([]model.ModelEnityResponse, 0)
	}

	helper.MessageOK = "model list retrieved"
	return response
}
