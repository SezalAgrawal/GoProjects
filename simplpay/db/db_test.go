package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testDB     Provider
	testDBURL  string
	testDBName string
)

func setup() {
	//testDBURL, testDBName = os.Getenv("TEST_DATABASE_URL"), os.Getenv("TEST_DATABASE_NAME")
	testDBURL, testDBName = "mongodb://localhost:27017", "simpltest"
	testDB = NewProvider(testDBURL, testDBName)
	if err := mustDeleteDB(testDBName, testDB.MongoClient); err != nil {
		panic(fmt.Sprintf("Unable to delete database, reason: %s", err))
	}
	testDB.MustCreateIndexes()
}

// teardown call this function only once after all the test are finished executing.
func teardown() {
	if err := mustDeleteDB(testDBName, testDB.MongoClient); err != nil {
		panic(fmt.Sprintf("Unable to delete database, reason: %s", err))
	}
	if err := testDB.Close(); err != nil {
		panic(fmt.Sprintf("Failed to close test database instance reason: %s", err))
	}
}

// This is a very dangerous func.
// Use it wisely and only in test files.
func mustDeleteDB(name string, fn func() (*mongo.Client, error)) error {
	c, err := fn()
	if err != nil {
		panic(err)
	}
	return c.Database(name).Drop(context.Background())
}

func TestMain(m *testing.M) {
	setup()
	ok := m.Run()
	teardown()
	os.Exit(ok)
}

func teardownFn(t *testing.T, dbName string, before func(dbName string, t *testing.T)) {
	before(dbName, t)
}

// beforeEachDBClean will clean the given database; can be used before executing every test.
var beforeEachDBClean = func(dbName string, t *testing.T) {
	mustDeleteDB(dbName, testDB.MongoClient)
	testDB.MustCreateIndexes()
}

// // cleanCollection will clean the given collection from database; can be used before executing every test.
// func cleanCollection(t *testing.T, dbName, colletionName string) {
// 	c, err := testDB.MongoClient()
// 	require.Nil(t, err)
// 	c.Database(dbName).Collection(colletionName).Drop(context.Background())
// 	testDB.MustCreateIndexes()
// 	// #TODO: can only delete records instead of dropping the collection
// }
