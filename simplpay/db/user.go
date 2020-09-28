package db

import (
	"context"
	"time"

	"github.com/goProjects/simplpay/errorops"
	"github.com/goProjects/simplpay/simplpay"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	usersCollection = "users"
)

func (db *database) CreateUser(ctx context.Context, u *simplpay.User) *errorops.Error {
	collection := db.session.Collection(usersCollection)
	context, cancelFunc := context.WithTimeout(ctx, ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	if _, err := collection.InsertOne(context, u); err != nil {
		return errorops.GetDBErr(err)
	}
	return nil
}

func (db *database) GetUser(ctx context.Context, name string) (*simplpay.User, *errorops.Error) {
	collection := db.session.Collection(usersCollection)
	context, cancelFunc := context.WithTimeout(ctx, ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	filter := bson.M{
		"name": name,
	}
	var u *simplpay.User
	if err := collection.FindOne(context, filter).Decode(&u); err != nil {
		return nil, errorops.GetDBErr(err)
	}
	return u, nil
}

// MustEnsureIndexesOnUsersCollection will panic if creating an indexes on target collection fail.
func (db *database) MustEnsureIndexesOnUsersCollection() {
	indexModels := make([]mongo.IndexModel, 0)
	indexModels = append(indexModels, getUniqueNameIndex())
	for _, indexModel := range indexModels {
		log.Infof("creating index on keys '%s' in %s collection", *indexModel.Options.Name, usersCollection)
		mustCreateIndex(indexModel, db.session.Collection(usersCollection))
	}
}

func getUniqueNameIndex() mongo.IndexModel {
	indexName := "name"
	unique := true
	keys := bsonx.Doc{
		{
			Key:   indexName,
			Value: bsonx.Int32(int32(asc)),
		},
	}
	index := mongo.IndexModel{}
	index.Keys = keys
	index.Options = &options.IndexOptions{
		Name:   &indexName,
		Unique: &unique,
	}
	return index
}
