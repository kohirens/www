package storage

import "github.com/kohirens/stdlib/logger"

// Storage Save data for long term.
type Storage interface {
	Load(filename string) ([]byte, error)
	Save(filename string, data []byte) error
	Remove(filename string) error
}

var Log = &logger.Standard{}
