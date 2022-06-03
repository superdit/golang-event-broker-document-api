package service

import "event-broker-document-api/model"

type BackendService interface {
	Insert(request model.BackendInsertRequest) (response model.BackendInsertResponse)
	Update(request model.BackendUpdateRequest) (response model.BackendUpdateResponse)
	Delete(request model.BackendDeleteRequest) (response model.EmptyResponse)
	List(query string) (response model.BackendListResponse)
}
