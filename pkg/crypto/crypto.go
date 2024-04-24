package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
)

const (
	tokenByteSize = 16
)

type Token string
type Hash []byte

func NewToken() (Token, error) {
	b := make([]byte, tokenByteSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return Token(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)), nil
}

func (t Token) ToHash() Hash {
	a := sha256.Sum256([]byte(t))
	return a[:]
}

func (h Hash) ToBase64() string {
	return base64.StdEncoding.EncodeToString(h)
}
