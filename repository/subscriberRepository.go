package repository

import (
	"event-broker-document-api/entity"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubscriberRepository interface {
	SessionTransaction() (*mongo.Client, *options.TransactionOptions)

	List(query string) (ents []entity.Subscriber)
	FindByBackendId(backendId string) (ents []entity.Subscriber)
	FindByEventId(eventId string) (ents []entity.Subscriber)
	FindByQueue(queue string) (ent entity.Subscriber)
	FindQueueExceptId(queue string, id string) (ent entity.Subscriber)
	FindById(id string) (ent entity.Subscriber)

	Insert(ent entity.Subscriber, contexts ...mongo.SessionContext)
	Update(ent entity.Subscriber, contexts ...mongo.SessionContext)
	Delete(entId string, contexts ...mongo.SessionContext)

	UpdateBackendByBackendId(entId string, entName string, contexts ...mongo.SessionContext)
	UpdateExchangeByEventId(entId string, entName string, contexts ...mongo.SessionContext)
}
