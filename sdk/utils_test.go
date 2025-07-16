package sdk

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
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create PEM encoded public key
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to marshal public key: %v", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	tests := []struct {
		name    string
		pemStr  string
		wantErr bool
	}{
		{
			name:    "valid PEM public key",
			pemStr:  string(publicKeyPEM),
			wantErr: false,
		},
		{
			name:    "invalid PEM format",
			pemStr:  "not a pem key",
			wantErr: true,
		},
		{
			name:    "empty string",
			pemStr:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publicKey, err := parsePublicKeyFromPEM(tt.pemStr)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if publicKey == nil {
				t.Error("Expected non-nil public key")
			}
		})
	}
}

func TestEncryptPKCS1v15(t *testing.T) {
	// Generate a test RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	publicKey := &privateKey.PublicKey
	testData := []byte("test encryption data")

	ciphertext, err := encryptPKCS1v15(publicKey, testData)
	if err != nil {
		t.Errorf("Encryption failed: %v", err)
		return
	}

	if len(ciphertext) == 0 {
		t.Error("Expected non-empty ciphertext")
	}

	// Test with data that's too large (should fail)
	largeData := make([]byte, 300) // Too large for 2048-bit RSA key with PKCS1v15 padding
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	_, err = encryptPKCS1v15(publicKey, largeData)
	if err == nil {
		t.Error("Expected error for data too large to encrypt")
	}
}
