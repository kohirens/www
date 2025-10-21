package gpg

import (
	"fmt"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// Capsule Encapsulate a GPG key.
//
//	This structure provides the application a simple way to use its own GPG
//	key to encrypt/decrypt data as needed. Normally the public and privateKey keys
//	would not be needed at the same time. The methods are simple wrappers to
//	their crypt counterparts in an effort aid human programmers make sense and
//	easily remember what they are doing. Error  s returned should be
//	comprehensible and indicate where in the program the problem lies.
type Capsule struct {
	PublicKey  *crypto.Key
	privateKey *crypto.Key
	pgp        *crypto.PGPHandle
}

// Decrypt armored encrypted message using the privateKey key and obtain the
// plaintext.
func (k *Capsule) Decrypt(encryptedMessage []byte) ([]byte, error) {
	decHandle, e1 := k.pgp.Decryption().DecryptionKey(k.privateKey).New()
	if e1 != nil {
		return nil, e1
	}

	// Clean up any traces of private key from memory.
	defer decHandle.ClearPrivateParams()

	decrypted, e2 := decHandle.Decrypt(encryptedMessage, crypto.Armor)
	if e2 != nil {
		return nil, e2
	}

	plaintext := decrypted.Bytes()

	return plaintext, nil
}

// Encrypt plaintext message using a public key
func (k *Capsule) Encrypt(message string) ([]byte, error) {
	encHandle, e1 := k.pgp.Encryption().Recipient(k.PublicKey).New()
	if e1 != nil {
		return nil, fmt.Errorf(stderr.Encryption, e1.Error())
	}

	pgpMessage, e2 := encHandle.Encrypt([]byte(message))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.Encrypt, e2.Error())
	}

	armored, e3 := pgpMessage.ArmorBytes()
	if e3 != nil {
		return nil, fmt.Errorf(stderr.ArmorBytes, e3.Error())
	}

	return armored, nil
}
