package storage

import "github.com/kohirens/stdlib/logger"

// Storage Save data for long term.
type Storage interface {
	// Exist Verification the file is in storage.
	Exist(name string) bool
	// Load Retrieve data from storage.
	Load(filename string) ([]byte, error)
	// Location Get the location in storage. This does not check for existence.
	Location(filename string) string
	// Save Write data to storage.
	Save(filename string, data []byte) error
	// Remove data from storage.
	Remove(filename string) error
}

var Log = &logger.Standard{}
