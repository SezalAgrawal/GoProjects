package db

import (
	"github.com/goProjects/simplpay/simplpay"
	"go.mongodb.org/mongo-driver/mongo"
)

// Provider holds all database methods.
type Provider interface {
	mongoDB
	simplpay.DatabaseProvider
	MustCreateIndexes()
	Close() error
	Ok() (bool, error)
	// Info() string
}

// DB is database type.
type database struct {
	session *mongo.Database
	client  *mongo.Client
	url     string
	mainDB  string
}

// NewProvider panics if database client fail to connect with given database server.
func NewProvider(dbURL, mainDB string) Provider {
	var err error
	db := &database{
		url:    dbURL,
		mainDB: mainDB,
	}
	if db.client, err = connectMongoDB(dbURL); err != nil {
		panic(err.Error())
	}
	db.session = db.client.Database(mainDB)
	return db
}

// MustCreateIndexes applies indexes on various collection at startup.
// It panics if creating indexes cause an error.
func (db *database) MustCreateIndexes() {
	db.MustEnsureIndexesOnUsersCollection()
	db.MustEnsureIndexesOnMerchantsCollection()
	db.MustEnsureIndexesOnTransactionsCollection()
}
