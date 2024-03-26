package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

type TokenSource struct {
}

func NewTokenSource() TokenSource {
	return TokenSource{}
}

func (s TokenSource) CreateApplicationToken() (string, error) {
	cryptoRand := make([]byte, 128) //nolint:gomnd
	_, err := rand.Read(cryptoRand)
	if err != nil {
		return "", errors.WithMessage(err, "crypto/rand read")
	}
	random := hex.EncodeToString(cryptoRand)

	return random, nil
}
