package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// ConnectionTimeoutInSecond is number of seconds after which the database operation will get timed out.
	ConnectionTimeoutInSecond = 200
)

// mongoDB interface provide methods which provide mongodb instance such as mongo client instance.
type mongoDB interface {
	// MongoClient ...
	MongoClient() (*mongo.Client, error)
}

func connectMongoDB(connectionString string) (*mongo.Client, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	return mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
}

// Ok check if database is running or not
func (db *database) MongoClient() (*mongo.Client, error) {
	if db.client == nil {
		return nil, errors.New("Mongodb client is not initialized")
	}
	return db.client, nil
}

// Ok check if database is running or not
func (db *database) Ok() (bool, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	if err := db.client.Ping(ctx, readpref.Primary()); err != nil {
		return false, err
	}
	return true, nil
}

// Close closes the databases connections
func (db *database) Close() error {
	return db.client.Disconnect(context.Background())
}

// IndexType is the type of index ASC or DESC
type IndexType int32

const (
	asc  IndexType = 1
	desc IndexType = -1
)

// mustCreateIndex will panic if creating an index on given collection fail
func mustCreateIndex(index mongo.IndexModel, c *mongo.Collection) {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	ctx, cancelFunc := context.WithTimeout(context.Background(), ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	if _, err := c.Indexes().CreateOne(ctx, index, opts); err != nil {
		panic(fmt.Sprintf("error while applying index to collection[%s], error[%s]", c.Name(), err.Error()))
	}
}
