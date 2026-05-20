package bsdkv3

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func parsePubkeyFromPEM(keyPEMString string) (*rsa.PublicKey, error) {
	keyPEMBytes := []byte(keyPEMString)

	block, _ := pem.Decode(keyPEMBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	if publicKey, ok := pub.(*rsa.PublicKey); ok {
		return publicKey, nil
	}

	return nil, errors.New("parsed key is not an RSA public key")
}

func encryptPKCS1v15(pubkey rsa.PublicKey, plainBytes []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &pubkey, plainBytes)

	if err != nil {
		return nil, fmt.Errorf("rsa encryption failed: %w", err)
	}
	return ciphertext, nil
}
