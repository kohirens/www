package backend

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/www/storage"
)

type Device struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
}

type Account struct {
	ID       string             `json:"id"`
	AppleID  string             `json:"apple_id"`
	Email    string             `json:"email"`
	FistName string             `json:"fist_name"`
	LastName string             `json:"last_name"`
	GoogleId string             `json:"google_id"`
	Devices  map[string]*Device `json:"devices"`
}

type AccountManager interface {
	Add(providerID, deviceID, providerName string) (*Account, error)
	Lookup(id string) (*Account, error)
}

type AccountExec struct {
	store storage.Storage
}

func Add() (*Account, error) {
	//TODO: generate a guid
	return &Account{
		ID:      "",
		Devices: make(map[string]*Device),
	}, nil
}

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
