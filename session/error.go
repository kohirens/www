package session

import (
	"fmt"
	"time"
)

type ExpiredError struct {
	exp time.Time
}

type StorageError struct{}

func (e StorageError) Error() string {
	return stderr.NoStorage
}

func (e ExpiredError) Error() string {
	return fmt.Sprintf(stderr.ExpiredCookie, e.exp.UTC().Format(time.RFC3339))
}

type NoSessionCookieError struct{}

func (e NoSessionCookieError) Error() string {
	return stderr.NoIDCookieFound
}

type RestoreError struct {
	msg string
}

func (e RestoreError) Error() string {
	return stderr.RestoreSession + e.msg
}

type InvalidIDError struct {
	id string
}

func (e InvalidIDError) Error() string {
	return stderr.InvalidSessionID + e.id
}
