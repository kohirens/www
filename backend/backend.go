package backend

import "github.com/kohirens/www/storage"

const (
	KeyAccountManager = "am"
	PrefixAccounts    = "accounts"
	PrefixGPGKey      = "secrets"
)

func NewAccountExec(store storage.Storage) *AccountExec {
	return &AccountExec{store: store}
}
