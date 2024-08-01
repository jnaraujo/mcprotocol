package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthDigest(t *testing.T) {
	assert.Equal(t, "4ed1f46bbe04bc756bcb17c0c7ce3e4632f06a48", AuthDigest("Notch"))
	assert.Equal(t, "-7c9d5b0044c130109a5d7b5fb5c317c02b4e28c1", AuthDigest("jeb_"))
	assert.Equal(t, "88e16a1019277b15d58faf0541e11910eb756f6", AuthDigest("simon"))
}
