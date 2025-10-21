package gpg

import (
	"fmt"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"os"
)

// LoadPublicKey Load a public GPG key from a file in armored format.
func LoadPublicKey(filename string) (*crypto.Key, error) {
	pubKeyData, e1 := os.ReadFile(filename)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ReadPublicKeyFile, e1.Error())
	}

	publicKey, e2 := crypto.NewKeyFromArmored(string(pubKeyData))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.ReadPublicKeyArmored, e2.Error())
	}

	return publicKey, nil
}

// NewCapsule Load a public and privateKey GPG key-pair, from files in
// armored format.
func NewCapsule(publicKeyFile, privateKeyFile, passphrase string) (*Capsule, error) {
	publicKey, e1 := LoadPublicKey(publicKeyFile)
	if e1 != nil {
		return nil, e1
	}
	privateKey, e2 := LoadPrivateKey(privateKeyFile, passphrase)
	if e2 != nil {
		return nil, e2
	}

	return &Capsule{
		pgp:        crypto.PGP(),
		PublicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

// LoadPrivateKey - load a privateKey key from filename and supply the passphrase
// of the privateKey key.
func LoadPrivateKey(filename, passphrase string) (*crypto.Key, error) {
	privateKeyData, e1 := os.ReadFile(filename)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ReadPrivateKeyFile, e1.Error())
	}

	privateKey, e2 := crypto.NewPrivateKeyFromArmored(
		string(privateKeyData),
		[]byte(passphrase),
	)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.ReadPrivateKeyArmored, e2.Error())
	}

	return privateKey, nil
}
