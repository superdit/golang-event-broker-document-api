package repository

import (
	"event-broker-document-api/entity"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository interface {
	SessionTransaction() (*mongo.Client, *options.TransactionOptions)

	List(query string) (ents []entity.Event)
	FindByBackendId(backendId string) (ent []entity.Event)
	FindByExchange(exchange string) (ent entity.Event)
	FindExchangeExceptId(exchange string, id string) (ent entity.Event)
	FindById(id string) (ent entity.Event)

	Insert(ent entity.Event, contexts ...mongo.SessionContext)
	Update(ent entity.Event, contexts ...mongo.SessionContext)
	Delete(entId string, contexts ...mongo.SessionContext)

	UpdateBackendByBackendId(entId string, entName string, contexts ...mongo.SessionContext)
	UpdateModelByModelId(entId string, entName string, entUrl string, contexts ...mongo.SessionContext)
}
