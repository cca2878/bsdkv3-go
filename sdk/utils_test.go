package bsdkv3

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestParsePublicKeyFromPEM(t *testing.T) {
	// Generate a test RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	// Convert public key to PEM format
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to marshal public key: %v", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	tests := []struct {
		name      string
		keyPEM    string
		wantError bool
	}{
		{
			name:      "valid PEM key",
			keyPEM:    string(publicKeyPEM),
			wantError: false,
		},
		{
			name:      "empty string",
			keyPEM:    "",
			wantError: true,
		},
		{
			name:      "invalid PEM format",
			keyPEM:    "not a pem key",
			wantError: true,
		},
		{
			name: "valid PEM but not a public key",
			keyPEM: `-----BEGIN CERTIFICATE-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA
-----END CERTIFICATE-----`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := parsePublicKeyFromPEM(tt.keyPEM)
			if (err != nil) != tt.wantError {
				t.Errorf("parsePublicKeyFromPEM() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && key == nil {
				t.Error("parsePublicKeyFromPEM() returned nil key without error")
			}
		})
	}
}

func TestEncryptPKCS1v15(t *testing.T) {
	// Generate a test RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	publicKey := &privateKey.PublicKey

	tests := []struct {
		name      string
		plaintext []byte
		wantError bool
	}{
		{
			name:      "valid plaintext",
			plaintext: []byte("test message"),
			wantError: false,
		},
		{
			name:      "empty plaintext",
			plaintext: []byte(""),
			wantError: false,
		},
		{
			name:      "large plaintext (should still work for 2048-bit key)",
			plaintext: make([]byte, 100),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := encryptPKCS1v15(publicKey, tt.plaintext)
			if (err != nil) != tt.wantError {
				t.Errorf("encryptPKCS1v15() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && ciphertext == nil {
				t.Error("encryptPKCS1v15() returned nil ciphertext without error")
			}
			if !tt.wantError && len(ciphertext) == 0 {
				t.Error("encryptPKCS1v15() returned empty ciphertext without error")
			}
		})
	}
}

func TestEncryptPKCS1v15_NilKey(t *testing.T) {
	_, err := encryptPKCS1v15(nil, []byte("test"))
	if err == nil {
		t.Error("encryptPKCS1v15() with nil key should return error")
	}
}
