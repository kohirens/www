package backend

import (
	"github.com/kohirens/www/storage"
	"testing"
)

func TestAccountExec_Lookup(t *testing.T) {
	fixedStore, _ := storage.NewLocalStorage(fixtureDir)

	tests := []struct {
		name    string
		store   storage.Storage
		id      string
		wantErr bool
	}{
		{
			"pull_account",
			fixedStore,
			"1234",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountExec{
				store: tt.store,
			}
			got, err := am.Lookup(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID != tt.id {
				t.Errorf("Lookup() got = %v, want %v", got, tt.id)
				return
			}
		})
	}
}

func TestAccountExec_Add(t *testing.T) {
	fixedStore, _ := storage.NewLocalStorage(fixtureDir)

	tests := []struct {
		name,
		providerID,
		providerName string
		store   storage.Storage
		device  *Device
		wantErr bool
	}{
		{
			"add_account",
			"1234",
			"google",
			fixedStore,
			&Device{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountExec{
				store: tt.store,
			}
			got, err := am.Add(tt.providerID, tt.providerName, tt.device)

			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.GoogleId != tt.providerID {
				t.Errorf("Add() got = %v, want %v", got.GoogleId, tt.providerID)
				return
			}
		})
	}
}
