package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
)

// code from https://gist.github.com/toqueteos/5372776

func AuthDigest(value string) string {
	hash := sha1.New()
	io.WriteString(hash, value)
	hashedValue := hash.Sum(nil)

	// Check for negative hashes
	negative := (hashedValue[0] & 0x80) == 0x80
	if negative {
		hashedValue = twosComplement(hashedValue)
	}

	res := strings.TrimLeft(hex.EncodeToString(hashedValue), "0")
	if negative {
		res = "-" + res
	}

	return res
}

func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}
