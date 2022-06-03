package service

import "event-broker-document-api/model"

type EventService interface {
	Insert(request model.EventInsertRequest) (response model.EventInsertResponse)
	Update(request model.EventUpdateRequest) (response model.EventUpdateResponse)
	Delete(request model.EventDeleteRequest) (response model.EmptyResponse)
	List(query string) (response model.EventListResponse)
}
