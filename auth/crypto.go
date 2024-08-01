package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type Crypto struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	privateKeyBytes []byte
	publicKeyBytes  []byte
}

func NewCrypto() (*Crypto, error) {
	crypto := new(Crypto)

	privateKey, publicKey, err := generateKeys()
	if err != nil {
		return nil, err
	}

	crypto.privateKey = privateKey
	crypto.publicKey = publicKey

	crypto.privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	crypto.publicKeyBytes, err = x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return crypto, nil
}

func generateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey

	privateKey.Precompute()
	err = privateKey.Validate()
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func (c *Crypto) PublicKeyBytes() []byte {
	return c.publicKeyBytes
}

func (c *Crypto) PrivateKeyBytes() []byte {
	return c.privateKeyBytes
}

func (c *Crypto) Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, data)
}

func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, c.privateKey, data)
}
