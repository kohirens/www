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
