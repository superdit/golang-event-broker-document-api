package repository

import (
	"event-broker-document-api/config"
	"event-broker-document-api/entity"
	"event-broker-document-api/exception"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type modelRepositoryImpl struct {
	Collection *mongo.Collection
	Database   mongo.Database
}

func NewModelRepository(database *mongo.Database) ModelRepository {
	return &modelRepositoryImpl{
		Collection: database.Collection("models"),
		Database:   *database,
	}
}

func (repository *modelRepositoryImpl) List(query string) (ents []entity.Model) {
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
		var entBson bson.M
		ent := entity.Model{}

		err = cursor.Decode(&entBson)
		exception.PanicIfNeeded(err)

		bsonBytes, err := bson.Marshal(entBson)
		exception.PanicIfNeeded(err)

		bson.Unmarshal(bsonBytes, &ent)
		ents = append(ents, ent)
	}
	return ents
}

func (repository *modelRepositoryImpl) FindByName(name string) (ent entity.Model) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"name": name}).Decode(&entBson)
	if err != nil {
		return entity.Model{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *modelRepositoryImpl) FindById(id string) (ent entity.Model) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	var entBson bson.M

	err := repository.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entBson)
	if err != nil {
		return entity.Model{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)

	return ent
}

func (repository *modelRepositoryImpl) FindNameExceptId(name string, id string) (ent entity.Model) {
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
		return entity.Model{}
	}

	bsonBytes, _ := bson.Marshal(entBson)
	bson.Unmarshal(bsonBytes, &ent)
	return ent
}

func (repository *modelRepositoryImpl) Insert(ent entity.Model) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	entBson, err := repository.Collection.InsertOne(ctx, ent)
	exception.PanicIfNeeded(err)

	bsonBytes, err := bson.Marshal(entBson)
	exception.PanicIfNeeded(err)

	bson.Unmarshal(bsonBytes, &ent)
}

func (repository *modelRepositoryImpl) Update(ent entity.Model) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	update := bson.M{"$set": ent}

	_, err := repository.Collection.UpdateOne(ctx, bson.M{"_id": ent.Id}, update)
	exception.PanicIfNeeded(err)
}

func (repository *modelRepositoryImpl) Delete(entId string) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	_, err := repository.Collection.DeleteOne(ctx, bson.M{"_id": entId})
	exception.PanicIfNeeded(err)
}
