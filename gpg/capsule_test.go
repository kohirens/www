package gpg

import (
	"testing"
)

func TestCapsule(t *testing.T) {
	tests := []struct {
		name        string
		publicFile  string
		privateFile string
		passphrase  string
		fingerprint string
		want        []byte
		wantErr     bool
	}{
		{
			"good",
			fixturesDir + "/gpg-test.public.asc",
			fixturesDir + "/gpg-test.private.asc",
			"test1234",
			"a353e4ecdb14ece84ad0fd909efb96ce70c116c5",
			[]byte("Salam"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e1 := NewCapsule(tt.publicFile, tt.privateFile, tt.passphrase)
			if (e1 != nil) != tt.wantErr {
				t.Errorf("LoadPrivateKey() error = %v, wantErr %v", e1, tt.wantErr)
				return
			}

			gotEncryptedMessage, e2 := got.Encrypt(tt.want)
			if (e2 != nil) != tt.wantErr {
				t.Errorf("Capsule.Encrypt() error = %v, wantErr %v", e2, tt.wantErr)
				return
			}

			gotMessage, e3 := got.Decrypt(gotEncryptedMessage)
			if (e3 != nil) != tt.wantErr {
				t.Errorf("Capsule.Decrypt() error = %v, wantErr %v", e3, tt.wantErr)
				return
			}

			if string(gotMessage) != string(tt.want) {
				t.Errorf("Capsule encrypt then decrypt returned %v, want %v", gotMessage, tt.want)
				return
			}
		})
	}
}
