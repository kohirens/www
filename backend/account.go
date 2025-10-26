package backend

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kohirens/www/storage"
)

type Account struct {
	AppleID  string             `json:"apple_id"`
	Devices  map[string]*Device `json:"devices"`
	Email    string             `json:"email"`
	FistName string             `json:"fist_name"`
	GoogleId string             `json:"google_id"`
	ID       string             `json:"id"`
	LastName string             `json:"last_name"`
}

type AccountManager interface {
	Add(providerID, providerName string, device *Device) (*Account, error)
	Lookup(id string) (*Account, error)
}

type AccountExec struct {
	store storage.Storage
}

// Add Make a new account.
func (am *AccountExec) Add(providerID, providerName string, device *Device) (*Account, error) {
	d := make(map[string]*Device)
	d[device.ID] = device

	// generate an account ID.
	id, e1 := uuid.NewV7()
	if e1 != nil {
		return nil, fmt.Errorf(stderr.UUID, e1.Error())
	}

	account := &Account{
		ID:      id.String(), //TODO: generate a guid
		Devices: d,
	}

	switch providerName {
	case "apple":
		account.AppleID = providerID
	case "google":
		account.GoogleId = providerID
	}

	accountBytes, e1 := json.Marshal(account)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, e1.Error())
	}

	if e := am.store.Save(account.ID, accountBytes); e != nil {
		return nil, e
	}

	return account, nil
}

// Lookup Search for an account in storage.
func (am *AccountExec) Lookup(id string) (*Account, error) {
	filename := KeyAccountPrefix + "/" + id + ".json"

	aData, e1 := am.store.Load(filename)
	if e1 != nil {
		return nil, &AccountNotFoundError{id}
	}

	account := &Account{}
	if e2 := json.Unmarshal(aData, &account); e2 != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, e2.Error())
	}

	return account, nil
}
