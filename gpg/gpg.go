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

	return PublicKey(string(pubKeyData))
}

// PublicKey Load a public GPG key from data in armored format.
func PublicKey(pubKeyData string) (*crypto.Key, error) {
	publicKey, e1 := crypto.NewKeyFromArmored(pubKeyData)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ReadPublicKeyArmored, e1.Error())
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

// NewCapsuleBytes Load GPG key from a data in armored format.
func NewCapsuleBytes(publicKeyData, privateKeyData []byte, passphrase string) (*Capsule, error) {
	publicKey, e1 := PublicKey(string(publicKeyData))
	if e1 != nil {
		return nil, e1
	}
	privateKey, e2 := PrivateKey(string(privateKeyData), passphrase)
	if e2 != nil {
		return nil, e2
	}

	return &Capsule{
		pgp:        crypto.PGP(),
		PublicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

// NewCapsuleString Load GPG key from a data in armored format.
func NewCapsuleString(publicKeyData, privateKeyData, passphrase string) (*Capsule, error) {
	publicKey, e1 := PublicKey(publicKeyData)
	if e1 != nil {
		return nil, e1
	}
	privateKey, e2 := PrivateKey(privateKeyData, passphrase)
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

	return PrivateKey(string(privateKeyData), passphrase)
}

// PrivateKey - load a privateKey key from []byte and the passphrase.
func PrivateKey(privateKeyData, passphrase string) (*crypto.Key, error) {
	privateKey, e1 := crypto.NewPrivateKeyFromArmored(
		privateKeyData,
		[]byte(passphrase),
	)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ReadPrivateKeyArmored, e1.Error())
	}

	return privateKey, nil
}
