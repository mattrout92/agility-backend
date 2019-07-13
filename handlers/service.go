package handlers

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ErrCouldNotGetNextQuestion is returned if the system could not get next question
var ErrCouldNotGetNextQuestion = errors.New("could not get next question id")

// ErrBadRequest is returned when the data provided does not allow for a first hit
type ErrBadRequest struct {
	err error
}

func (e ErrBadRequest) Error() string {
	return fmt.Sprintf("could not format a session based on input data: [%v]", e.err)
}

// Storer represents the methods to access datastore
type Storer interface {
	GetQuery(queryName string) string
	SQLX() *sqlx.DB
}

// Service contains a storer
type Service struct {
	Store Storer
}
