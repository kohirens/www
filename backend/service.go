package backend

import (
	"fmt"
)

func NewServiceManager() ServiceManager {
	return &Servicer{
		services: make(map[string]any),
	}
}

// ServiceManager handles storing and retrieval of dependencies when an endpoint
// handler function is called. Granting access to objects such as Session,
// Database, and whatever else may be useful.
type ServiceManager interface {
	Add(name string, service any)
	Get(name string) (any, error)
}

// Key acts as a way to provide information on the type of structure we are
// storing. This is done because Go does NOT type parameters on methods, just
// functions. You'll see later that we wrap methods for storing and retrieving
// services from the service manager.
type Key[T any] struct {
	name string
}

// NewKey associates each key in the service manager list with the type that it
// stores.
func NewKey[T any](name string) Key[T] {
	return Key[T]{name}
}

// Servicer the default service manager provided by the backend package.
type Servicer struct {
	services map[string]any
}

// Add adds a service to the list of services.
//
//	NOTE: This is for internal use only, it may be made private in the future.
//	Please use the Store[T any]() function.
func (m *Servicer) Add(name string, service any) {
	m.services[name] = service
}

// Get Retrieves a service from the list of services.
//
//	NOTE: This is for internal use only, it may be made private in the future.
//	Please use the Retrieve[T any]() function.
func (m *Servicer) Get(name string) (any, error) {
	x, ok := m.services[name]
	if !ok {
		return nil, fmt.Errorf(stderr.ServiceNotFound, name)
	}

	return x, nil
}

// Store adds a service to a map for access from an endpoint function.
func Store[T any](m ServiceManager, key Key[T], service T) {
	m.Add(key.name, service)
}

// Retrieve Retrieves a service from the list of services.
func Retrieve[T any](m ServiceManager, key Key[T]) (T, error) {
	stored, e1 := m.Get(key.name)
	if e1 != nil {
		var nilVar T
		return nilVar, e1
	}

	service, ok := stored.(T)
	if !ok {
		var nilVar T
		return nilVar, fmt.Errorf("%v", stderr.ServicePointer)
	}

	return service, nil
}
