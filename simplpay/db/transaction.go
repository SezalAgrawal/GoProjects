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
	transactionsCollection = "transactions"
)

func (db *database) CreateTransaction(ctx context.Context, t *simplpay.Transaction, userOwnedAmount float64, updatedOn time.Time) *errorops.Error {
	// transaction starts here.
	err := db.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		// start transaction.
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		// update user owned amount.
		selector := bson.M{
			"name": t.UserName,
		}
		updateFields := bson.M{
			"$set": bson.M{
				"owned_amount": userOwnedAmount,
				"updated_on":   updatedOn,
			},
		}
		if result := db.session.Collection(usersCollection).FindOneAndUpdate(sessionContext, selector, updateFields); result.Err() != nil {
			return result.Err()
		}

		// create transaction.
		if _, err := db.session.Collection(transactionsCollection).InsertOne(sessionContext, t); err != nil {
			return err
		}

		// commit transaction if no error occurs.
		if err := sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	},
	)
	// transaction ends here.

	if err != nil {
		return errorops.GetDBErr(err)
	}
	return nil
}

// MustEnsureIndexesOnTransactionsCollection will panic if creating an indexes on target collection fail.
func (db *database) MustEnsureIndexesOnTransactionsCollection() {
	indexModels := make([]mongo.IndexModel, 0)
	indexModels = append(indexModels, getUniqueIDIndex())
	for _, indexModel := range indexModels {
		log.Infof("creating index on keys '%s' in %s collection", *indexModel.Options.Name, transactionsCollection)
		mustCreateIndex(indexModel, db.session.Collection(transactionsCollection))
	}
}

func getUniqueIDIndex() mongo.IndexModel {
	indexName := "id"
	unique := true
	keys := bsonx.Doc{
		{
			Key:   indexName,
			Value: bsonx.Int32(int32(desc)),
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
