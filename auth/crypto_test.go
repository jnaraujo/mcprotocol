package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	crypto, err := NewCrypto()
	assert.Nil(t, err)

	err = crypto.privateKey.Validate()
	assert.Nil(t, err)

	assert.Equal(t, 1024, crypto.privateKey.N.BitLen())
	assert.Equal(t, crypto.privateKey.PublicKey, *crypto.publicKey)
}
