package repository

import (
	"event-broker-document-api/entity"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BackendRepository interface {
	SessionTransaction() (*mongo.Client, *options.TransactionOptions)

	List(query string) (ents []entity.Backend)
	FindByName(name string) (ent entity.Backend)
	FindNameExceptId(name string, id string) (ent entity.Backend)
	FindById(id string) (ent entity.Backend)
	Insert(ent entity.Backend, contexts ...mongo.SessionContext)
	Update(ent entity.Backend, contexts ...mongo.SessionContext)
	Delete(entId string, contexts ...mongo.SessionContext)
}
