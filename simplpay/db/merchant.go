package db

import (
	"context"
	"time"

	"github.com/goProjects/simplpay/errorops"
	"github.com/goProjects/simplpay/simplpay"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	merchantsCollection = "merchants"
)

func (db *database) CreateMerchant(ctx context.Context, m *simplpay.Merchant) *errorops.Error {
	collection := db.session.Collection(merchantsCollection)
	context, cancelFunc := context.WithTimeout(ctx, ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	if _, err := collection.InsertOne(context, m); err != nil {
		return errorops.GetDBErr(err)
	}
	return nil
}

func (db *database) GetMerchant(ctx context.Context, name string) (*simplpay.Merchant, *errorops.Error) {
	collection := db.session.Collection(merchantsCollection)
	context, cancelFunc := context.WithTimeout(ctx, ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	filter := bson.M{
		"name": name,
	}
	var m *simplpay.Merchant
	if err := collection.FindOne(context, filter).Decode(&m); err != nil {
		return nil, errorops.GetDBErr(err)
	}
	return m, nil
}

// MustEnsureIndexesOnUsersCollection will panic if creating an indexes on target collection fail.
func (db *database) MustEnsureIndexesOnMerchantsCollection() {
	indexModels := make([]mongo.IndexModel, 0)
	indexModels = append(indexModels, getUniqueNameIndex())
	for _, indexModel := range indexModels {
		log.Infof("creating index on keys '%s' in %s collection", *indexModel.Options.Name, merchantsCollection)
		mustCreateIndex(indexModel, db.session.Collection(merchantsCollection))
	}
}
