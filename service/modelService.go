package service

import "event-broker-document-api/model"

type ModelService interface {
	Insert(request model.ModelInsertRequest) (response model.ModelInsertResponse)
	Update(request model.ModelUpdateRequest) (response model.ModelUpdateResponse)
	Delete(request model.ModelDeleteRequest) (response model.EmptyResponse)
	List(query string) (response model.ModelListResponse)
}
