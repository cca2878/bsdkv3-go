package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/cca2878/bsdkv3-go/internal/apierr"
)

func parsePubkeyFromPEM(keyPEMString string) (*rsa.PublicKey, error) {
	keyPEMBytes := []byte(keyPEMString)

	block, _ := pem.Decode(keyPEMBytes)
	if block == nil {
		return nil, fmt.Errorf("%w: PEM 解码失败", apierr.ErrCipher)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: 解析公钥失败: %w", apierr.ErrCipher, err)
	}

	if publicKey, ok := pub.(*rsa.PublicKey); ok {
		return publicKey, nil
	}

	return nil, fmt.Errorf("%w: 不是 RSA 公钥", apierr.ErrCipher)
}

func encryptPKCS1v15(pubkey rsa.PublicKey, plainBytes []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &pubkey, plainBytes)

	if err != nil {
		return nil, fmt.Errorf("%w: RSA 加密失败: %w", apierr.ErrCipher, err)
	}
	return ciphertext, nil
}
