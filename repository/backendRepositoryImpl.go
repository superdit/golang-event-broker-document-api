package repository

import (
	"context"
	"errors"
	"event-broker-document-api/config"
	"event-broker-document-api/entity"
	"event-broker-document-api/exception"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type backendRepositoryImpl struct {
	Collection         *mongo.Collection
	Database           mongo.Database
	DatabaseContext    context.Context
	DatabaseCancelFunc context.CancelFunc
}

func NewBackendRepository(database *mongo.Database) BackendRepository {
	return &backendRepositoryImpl{
		Collection:         database.Collection("backends"),
		Database:           *database,
		DatabaseContext:    nil,
		DatabaseCancelFunc: nil,
	}
}

func (repository *backendRepositoryImpl) useCustomContext(contexts []mongo.SessionContext) bool {
	ctxLen := len(contexts)

	if ctxLen == 1 {
		repository.DatabaseContext = contexts[0]
		return true
	}

	if ctxLen > 1 {
		err := errors.New("too many contexts")
		exception.PanicIfNeeded(err)
	}

	repository.DatabaseContext, repository.DatabaseCancelFunc = config.NewMongoContext()
	return false
}

func (repository *backendRepositoryImpl) SessionTransaction() (*mongo.Client, *options.TransactionOptions) {
	dbClient := repository.Database.Client()
	return dbClient, config.MongoTransOption()
}

func (repository *backendRepositoryImpl) List(query string) (ents []entity.Backend) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"name": bson.M{
				"$regex":   query,
				"$options": "i",
			},
		}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{Key: "name", Value: 1},
	})

	cursor, err := repository.Collection.Find(ctx, filter, findOptions)
	exception.PanicIfNeeded(err)

	for cursor.Next(ctx) {
		var entityBson bson.M
		entity := entity.Backend{}

		err = cursor.Decode(&entityBson)
		exception.PanicIfNeeded(err)

		bsonBytes, err := bson.Marshal(entityBson)
		exception.PanicIfNeeded(err)

		bson.Unmarshal(bsonBytes, &entity)
		ents = append(ents, entity)
	}
	return ents
}

func (repository *backendRepositoryImpl) FindByName(name string) (ent entity.Backend) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"name": name}).Decode(&entBson)
	if err != nil {
		return entity.Backend{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *backendRepositoryImpl) FindNameExceptId(name string, id string) (ent entity.Backend) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{
		"_id": bson.M{
			"$ne": id,
		},
		"name": bson.M{
			"$eq": name,
		},
	}

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, filter).Decode(&entBson)
	if err != nil {
		return entity.Backend{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)
	return ent
}

func (repository *backendRepositoryImpl) FindById(id string) (ent entity.Backend) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entBson)
	if err != nil {
		return entity.Backend{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *backendRepositoryImpl) Insert(ent entity.Backend, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	entBson, err := repository.Collection.InsertOne(repository.DatabaseContext, ent)
	exception.PanicIfNeeded(err)

	bsonBytes, err := bson.Marshal(entBson)
	exception.PanicIfNeeded(err)

	bson.Unmarshal(bsonBytes, &ent)
}

func (repository *backendRepositoryImpl) Update(ent entity.Backend, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}
	update := bson.M{"$set": ent}

	_, err := repository.Collection.UpdateOne(repository.DatabaseContext, bson.M{"_id": ent.Id}, update)
	exception.PanicIfNeeded(err)
}

func (repository *backendRepositoryImpl) Delete(entId string, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	_, err := repository.Collection.DeleteOne(repository.DatabaseContext, bson.M{"_id": entId})
	exception.PanicIfNeeded(err)
}
