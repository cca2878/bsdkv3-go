// Package bsdkv3 provides cryptographic utilities for the BSdk V3 client.
//
// This file contains utilities for:
//   - RSA public key parsing from PEM format
//   - PKCS#1 v1.5 encryption
//   - Password encryption with public key
package bsdkv3

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// Cipher

// parsePublicKeyFromPEM parses a PEM encoded public key from a string.
// It supports PKIX format which includes PKCS#1.
func parsePublicKeyFromPEM(keyPEMString string) (*rsa.PublicKey, error) {
	// Convert the string to bytes, as pem.Decode operates on bytes.
	keyPEMBytes := []byte(keyPEMString)

	block, _ := pem.Decode(keyPEMBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Assert the type to an RSA public key
	if publicKey, ok := pub.(*rsa.PublicKey); ok {
		return publicKey, nil
	}

	return nil, errors.New("parsed key is not an RSA public key")
}

// encryptPKCS1v15 encrypts data using RSA PKCS#1 v1.5 padding.
func encryptPKCS1v15(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("public key cannot be nil")
	}

	// rand.Reader is required for PKCS#1 v1.5 padding encryption
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("rsa encryption failed: %w", err)
	}
	return ciphertext, nil
}
