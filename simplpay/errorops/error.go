package errorops

import (
	"errors"
	"strconv"
	"strings"

	"go.caringcompany.co/simpliclaim/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// Error handles the error
type Error struct {
	Err  error
	Code string
}

const (
	// duplicateKeyError is error code which is thrown when you try
	// to insert duplicate field in an unique field index collection
	duplicateKeyError = 11000
)

// list of error constants
const (
	Unknown    string = "unknown"
	Processing string = "processing"
	NotFound   string = "not-found"
	Conflict   string = "conflict"
	Structure  string = "structure"
)

// dupErr returns whether err informs of a duplicate key error because
// a primary key index or a secondary unique index already has an entry
// with the given value.
func dupErr(err error) bool {
	if strings.Contains(err.Error(), strconv.Itoa(duplicateKeyError)) {
		return true
	}
	return false
}

func commandErr(err error) bool {
	ce := mongo.CommandError{}
	if errors.As(err, &ce) {
		return true
	}
	return false
}

func marshalErr(err error) bool {
	me := mongo.MarshalError{}
	if errors.As(err, &me) {
		return true
	}
	return false
}

func noDocumentErr(err error) bool {
	if err == mongo.ErrNoDocuments {
		return true
	}
	return false
}

// GetDBErr gets errorops Error object from DB error
func GetDBErr(err error) *Error {
	log.Errorw(err.Error(), log.ErrMsgKey, err)
	errCode := Unknown
	if noDocumentErr(err) {
		errCode = NotFound
	} else if dupErr(err) {
		errCode = Conflict
	} else if marshalErr(err) {
		errCode = Structure
	} else if commandErr(err) {
		errCode = Processing
	}

	e := &Error{
		Err:  err,
		Code: errCode,
	}
	return e
}
