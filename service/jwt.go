package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Jwt struct {
	secret string
}

func NewJwt(secret string) Jwt {
	return Jwt{
		secret: secret,
	}
}

func (s Jwt) CreateApplicationToken(appId int, expireTime int) (string, error) {
	random, err := s.generateSalt()
	if err != nil {
		return "", errors.WithMessage(err, "generate salt")
	}

	created := time.Now().UTC()
	claims := jwt.MapClaims{
		"appId": appId,
		"iat":   created.Unix(),
		"salt":  random,
	}
	if expireTime > 0 {
		claims["exp"] = created.Add(time.Millisecond * time.Duration(expireTime)).Unix()
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).
		SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.WithMessage(err, "signed token")
	}

	return token, nil
}

func (Jwt) generateSalt() (string, error) {
	cryptoRand := make([]byte, 16) //nolint:gomnd
	_, err := rand.Read(cryptoRand)
	if err != nil {
		return "", errors.WithMessage(err, "crypto/rand read")
	}
	random := hex.EncodeToString(cryptoRand)

	return random, nil
}
