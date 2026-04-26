package session

import (
	"errors"
	"fmt"
	"time"
)

type ExpiredError struct {
	exp time.Time
}

var StorageError = errors.New(stderr.NoStorage)

func (e ExpiredError) Error() string {
	return fmt.Sprintf(stderr.ExpiredCookie, e.exp.UTC().Format(time.RFC3339))
}

type NoSessionError struct {
}

func (e NoSessionError) Error() string {
	return "no session"
}

type RestoreError struct {
	msg string
}

func (e RestoreError) Error() string {
	return "failed to restore session " + e.msg
}

type InvalidIDError struct {
	id string
}

func (e InvalidIDError) Error() string {
	return "invalid session id " + e.id
}
