package service

import "event-broker-document-api/model"

type SubscriberService interface {
	Insert(request model.SubscriberInsertRequest) (response model.SubscriberInsertResponse)
	Update(request model.SubscriberUpdateRequest) (response model.SubscriberUpdateResponse)
	Delete(request model.SubscriberDeleteRequest) (response model.EmptyResponse)
	List(query string) (response model.SubscriberListResponse)
}
