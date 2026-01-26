// Package session works to extend HTTP State Management Mechanism beyond the
// HTTP cookie header storage. Such as using files on servers that have file
// storage access or database storage for those without. The latter options
// can add more latency, so please consider your options and use or build an
// implementation according to your use case. For clarity on the subject of
// HTTP State Management please review the RFC at
// https://datatracker.ietf.org/doc/html/rfc6265
package session

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kohirens/stdlib/logger"
)

type Data struct {
	Id         string    `json,bson:"session_id"`
	Expiration time.Time `json,bson:"expiration"`
	Items      Store     `json,bson:"session_data"`
}

// Storage An interface medium for storing the session data to anyplace an
// implementer see fit. An implementor should especially take into consideration
// sensitive data pertaining to the clients session. This simple interface does
// implementation of encryption for Save and decryption for Load. Use this
// to implement storage for mediums like File, Database, In-memory cache, etc.
type Storage interface {
	// Load The session from storage.
	// No matter the storage medium this should always return JSON as a byte array.
	Load(id string) ([]byte, error)

	// Save The session data to the storage medium.
	Save(id string, data []byte) error

	// Remove Delete data from storage.
	Remove(key string) error
}

// Store Model for short term storage in memory (not intended for long
// term storage).
type Store map[string][]byte

const (
	IDKey = "_sid_" // IDKey Cookie name.
)

var (
	// ExtendTime How much time the session is extended when a user loads a
	// page after the initial start of the session
	ExtendTime     = 5 * time.Minute
	Log            = logger.Standard{}
	IDCookiePath   = "/"     // IDCookiePath Any path in the domain.
	IDCookieDomain = ""      // IDCookieDomain Default to the entire domain.
	Suffix         = ".json" // Optional file extension to append to the session save file.
)

// GenerateID A unique session ID
//
//	Panics if an ID cannot be generated.
func GenerateID() string {
	id, e1 := uuid.NewV7()
	if e1 != nil {
		msg := fmt.Sprintf(stderr.UUID, e1.Error())
		panic(msg)
	}
	return id.String()
}

// NewManager Initialize a new session manager to handle session save, restore, get, and set.
func NewManager(storage Storage, location string, expiration time.Duration) *Manager {
	return &Manager{
		data: &Data{
			GenerateID(),
			time.Now().UTC().Add(expiration),
			make(Store, 100),
		},
		storage:    storage,
		hasUpdates: false,
		location:   location,
	}
}
