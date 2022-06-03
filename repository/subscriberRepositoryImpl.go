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

type subscriberRepositoryImpl struct {
	Collection         *mongo.Collection
	Database           mongo.Database
	DatabaseContext    context.Context
	DatabaseCancelFunc context.CancelFunc
}

func NewSubscriberRepository(database *mongo.Database) SubscriberRepository {
	return &subscriberRepositoryImpl{
		Collection:         database.Collection("subscribers"),
		Database:           *database,
		DatabaseContext:    nil,
		DatabaseCancelFunc: nil,
	}
}

func (repository *subscriberRepositoryImpl) useCustomContext(contexts []mongo.SessionContext) bool {
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

func (repository *subscriberRepositoryImpl) SessionTransaction() (*mongo.Client, *options.TransactionOptions) {
	dbClient := repository.Database.Client()
	return dbClient, config.MongoTransOption()
}

func (repository *subscriberRepositoryImpl) List(query string) (ents []entity.Subscriber) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": bson.A{
				bson.M{
					"queue": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
				bson.M{
					"_backend_subscriber_name": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
				bson.M{
					"_exchange": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
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
		entity := entity.Subscriber{}

		err = cursor.Decode(&entityBson)
		exception.PanicIfNeeded(err)

		bsonBytes, err := bson.Marshal(entityBson)
		exception.PanicIfNeeded(err)

		bson.Unmarshal(bsonBytes, &entity)
		ents = append(ents, entity)
	}
	return ents
}

func (repository *subscriberRepositoryImpl) FindByBackendId(backendId string) (ents []entity.Subscriber) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{"_backend_subscriber_id": backendId}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{Key: "_backend_subscriber_name", Value: 1},
	})

	cursor, err := repository.Collection.Find(ctx, filter, findOptions)
	exception.PanicIfNeeded(err)

	for cursor.Next(ctx) {
		var entityBson bson.M
		entity := entity.Subscriber{}

		err = cursor.Decode(&entityBson)
		exception.PanicIfNeeded(err)

		bsonBytes, err := bson.Marshal(entityBson)
		exception.PanicIfNeeded(err)

		bson.Unmarshal(bsonBytes, &entity)
		ents = append(ents, entity)
	}

	return ents
}

func (repository *subscriberRepositoryImpl) FindByEventId(eventId string) (ents []entity.Subscriber) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{"_event_id": eventId}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{Key: "_exchange", Value: 1},
	})

	cursor, err := repository.Collection.Find(ctx, filter, findOptions)
	exception.PanicIfNeeded(err)

	for cursor.Next(ctx) {
		var entityBson bson.M
		entity := entity.Subscriber{}

		err = cursor.Decode(&entityBson)
		exception.PanicIfNeeded(err)

		bsonBytes, err := bson.Marshal(entityBson)
		exception.PanicIfNeeded(err)

		bson.Unmarshal(bsonBytes, &entity)
		ents = append(ents, entity)
	}

	return ents
}

func (repository *subscriberRepositoryImpl) FindByQueue(queue string) (ent entity.Subscriber) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"queue": queue}).Decode(&entBson)
	if err != nil {
		return entity.Subscriber{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *subscriberRepositoryImpl) FindQueueExceptId(queue string, id string) (ent entity.Subscriber) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	filter := bson.M{
		"_id": bson.M{
			"$ne": id,
		},
		"queue": bson.M{
			"$eq": queue,
		},
	}

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, filter).Decode(&entBson)
	if err != nil {
		return entity.Subscriber{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)
	return ent
}

func (repository *subscriberRepositoryImpl) FindById(id string) (ent entity.Subscriber) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entBson)
	if err != nil {
		return entity.Subscriber{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *subscriberRepositoryImpl) Insert(ent entity.Subscriber, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	entBson, err := repository.Collection.InsertOne(repository.DatabaseContext, ent)
	exception.PanicIfNeeded(err)

	bsonBytes, err := bson.Marshal(entBson)
	exception.PanicIfNeeded(err)

	bson.Unmarshal(bsonBytes, &ent)
}

func (repository *subscriberRepositoryImpl) Update(ent entity.Subscriber, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	update := bson.M{"$set": ent}

	_, err := repository.Collection.UpdateOne(repository.DatabaseContext, bson.M{"_id": ent.Id}, update)
	exception.PanicIfNeeded(err)
}

func (repository *subscriberRepositoryImpl) Delete(entId string, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}
	_, err := repository.Collection.DeleteOne(repository.DatabaseContext, bson.M{"_id": entId})
	exception.PanicIfNeeded(err)
}

func (repository *subscriberRepositoryImpl) UpdateBackendByBackendId(entId string, entName string, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	update := bson.M{"$set": bson.M{
		"_backend_subscriber_name": entName,
	}}

	_, err := repository.Collection.UpdateMany(repository.DatabaseContext, bson.M{"_backend_subscriber_id": entId}, update)
	exception.PanicIfNeeded(err)
}

func (repository *subscriberRepositoryImpl) UpdateExchangeByEventId(entId string, entName string, contexts ...mongo.SessionContext) {
	if !repository.useCustomContext(contexts) {
		defer repository.DatabaseCancelFunc()
	}

	update := bson.M{"$set": bson.M{
		"_exchange": entName,
	}}

	_, err := repository.Collection.UpdateMany(repository.DatabaseContext, bson.M{"_event_id": entId}, update)
	exception.PanicIfNeeded(err)
}
