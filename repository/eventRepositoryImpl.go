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

type eventRepositoryImpl struct {
	Collection         *mongo.Collection
	Database           mongo.Database
	DatabaseContext    context.Context
	DatabaseCancelFunc context.CancelFunc
}

func NewEventRepository(database *mongo.Database) EventRepository {
	return &eventRepositoryImpl{
		Collection:         database.Collection("events"),
		Database:           *database,
		DatabaseContext:    nil,
		DatabaseCancelFunc: nil,
	}
}

func (repository *eventRepositoryImpl) useCustomContext(contexts []mongo.SessionContext) bool {
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

func (repository *eventRepositoryImpl) SessionTransaction() (*mongo.Client, *options.TransactionOptions) {
	dbClient := repository.Database.Client()
	return dbClient, config.MongoTransOption()
}

func (repository *eventRepositoryImpl) List(query string) (ents []entity.Event) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": bson.A{
				bson.M{
					"exchange": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
				bson.M{
					"_backend_publisher_name": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
				bson.M{
					"_model_name": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
			},
		}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{Key: "exchange", Value: 1},
	})

	cursor, err := repository.Collection.Find(ctx, filter, findOptions)
	exception.PanicIfNeeded(err)

	for cursor.Next(ctx) {
		var entityBson bson.M
		entity := entity.Event{}

		err = cursor.Decode(&entityBson)
		exception.PanicIfNeeded(err)

		bsonBytes, err := bson.Marshal(entityBson)
		exception.PanicIfNeeded(err)

		bson.Unmarshal(bsonBytes, &entity)
		ents = append(ents, entity)
	}
	return ents
}

func (repository *eventRepositoryImpl) FindByBackendId(backendId string) (ents []entity.Event) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{"_backend_publisher_id": backendId}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{Key: "exchange", Value: 1},
	})

	cursor, err := repository.Collection.Find(ctx, filter, findOptions)
	exception.PanicIfNeeded(err)

	for cursor.Next(ctx) {
		var entityBson bson.M
		entity := entity.Event{}

		err = cursor.Decode(&entityBson)
		exception.PanicIfNeeded(err)

		bsonBytes, err := bson.Marshal(entityBson)
		exception.PanicIfNeeded(err)

		bson.Unmarshal(bsonBytes, &entity)
		ents = append(ents, entity)
	}

	return ents
}

func (repository *eventRepositoryImpl) FindByExchange(exchange string) (ent entity.Event) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"exchange": exchange}).Decode(&entBson)
	if err != nil {
		return entity.Event{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *eventRepositoryImpl) FindExchangeExceptId(exchange string, id string) (ent entity.Event) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{
		"_id": bson.M{
			"$ne": id,
		},
		"exchange": bson.M{
			"$eq": exchange,
		},
	}

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, filter).Decode(&entBson)
	if err != nil {
		return entity.Event{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)
	return ent
}

func (repository *eventRepositoryImpl) FindById(id string) (ent entity.Event) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entBson)
	if err != nil {
		return entity.Event{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *eventRepositoryImpl) Insert(ent entity.Event, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	entBson, err := repository.Collection.InsertOne(repository.DatabaseContext, ent)
	exception.PanicIfNeeded(err)

	bsonBytes, err := bson.Marshal(entBson)
	exception.PanicIfNeeded(err)

	bson.Unmarshal(bsonBytes, &ent)
}

func (repository *eventRepositoryImpl) Update(ent entity.Event, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	update := bson.M{"$set": ent}

	_, err := repository.Collection.UpdateOne(repository.DatabaseContext, bson.M{"_id": ent.Id}, update)
	exception.PanicIfNeeded(err)
}

func (repository *eventRepositoryImpl) Delete(entId string, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}
	_, err := repository.Collection.DeleteOne(repository.DatabaseContext, bson.M{"_id": entId})
	exception.PanicIfNeeded(err)
}

func (repository *eventRepositoryImpl) UpdateBackendByBackendId(entId string, entName string, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	update := bson.M{"$set": bson.M{
		"_backend_publisher_name": entName,
	}}

	_, err := repository.Collection.UpdateMany(repository.DatabaseContext, bson.M{"_backend_publisher_id": entId}, update)
	exception.PanicIfNeeded(err)
}

func (repository *eventRepositoryImpl) UpdateModelByModelId(entId string, entName string, entUrl string, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	update := bson.M{"$set": bson.M{
		"_model_name": entName,
		"_model_url":  entUrl,
	}}

	_, err := repository.Collection.UpdateMany(repository.DatabaseContext, bson.M{"_model_id": entId}, update)
	exception.PanicIfNeeded(err)
}
