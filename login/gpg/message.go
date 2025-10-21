package gpg

var stderr = struct {
	ArmorBytes,
	Encrypt,
	Encryption,
	ReadPrivateKeyArmored,
	ReadPrivateKeyFile,
	ReadPublicKey,
	ReadPublicKeyFile,
	ReadPublicKeyArmored string
}{
	ArmorBytes:            "cannot get armored key data %v",
	Encrypt:               "cannot encrypt message %v",
	Encryption:            "",
	ReadPrivateKeyArmored: "reading privateKey key armored %v",
	ReadPrivateKeyFile:    "reading privateKey key file %v",
	ReadPublicKey:         "reading public key %v",
	ReadPublicKeyFile:     "reading public key file %v",
	ReadPublicKeyArmored:  "reading public key armored %v",
}
