package gpg

import (
	"testing"
)

const (
	fixturesDir = "testdata"
	tmpDir      = "tmp"
)

func TestLoadPublicKey(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		fingerprint string
		wantErr     bool
	}{
		{
			"key-not-found",
			"key-not-found.asc",
			"",
			true,
		},
		{
			"good",
			fixturesDir + "/gpg-test.public.asc",
			"a353e4ecdb14ece84ad0fd909efb96ce70c116c5",
			false,
		},
		{
			"tampered",
			fixturesDir + "/gpg-test-tampered.public.asc",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := LoadPublicKey(tt.filename)
			if (e != nil) != tt.wantErr {
				t.Errorf("LoadPublicKey() error = %v, wantErr %v", e, tt.wantErr)
				return
			}

			if got != nil && got.GetFingerprint() != tt.fingerprint {
				t.Errorf("LoadPublicKey() got fingerprint %v, want %v", got.GetFingerprint(), tt.fingerprint)
				return
			}
		})
	}
}

func TestLoadPrivateKey(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		passphrase  string
		fingerprint string
		wantErr     bool
	}{
		{
			"good",
			fixturesDir + "/gpg-test.private.asc",
			"test1234",
			"a353e4ecdb14ece84ad0fd909efb96ce70c116c5",
			false,
		},
		{
			"key-not-found",
			"key-not-found.asc",
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadPrivateKey(tt.filename, tt.passphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.GetFingerprint() != tt.fingerprint {
				t.Errorf("LoadPublicKey() got fingerprint %v, want %v", got.GetFingerprint(), tt.fingerprint)
				return
			}
		})
	}
}
