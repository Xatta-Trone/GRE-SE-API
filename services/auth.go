package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/o1egl/paseto"
)

type AuthInterface interface {
	GenerateToken(s string) (string, error)
	ValidateToken(s string) string
}

func GenerateToken(s string) (string, error) {

	// get the key
	key := os.Getenv("AUTH_KEY")

	if key == "" {
		panic("AUTH_KEY not found")
	}

	ttl := os.Getenv("AUTH_TTL")

	if ttl == "" {
		panic("AUTH_TTL not found")
	}

	ttlInt, err := strconv.Atoi(ttl)

	if err != nil {
		panic("Invalid AUTH_TTL integer")
	}

	fmt.Println(ttlInt)

	symmetricKey := []byte(key) // Must be 32 bytes
	now := time.Now()
	exp := now.Add(time.Duration(ttlInt) * time.Second)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   "gre-sentence-equivalence.com",
		Issuer:     "gre-sentence-equivalence.com",
		Jti:        s,
		Subject:    s,
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	// Add custom claim    to the token
	jsonToken.Set("data", "this is a signed message")
	footer := "gre-sentence-equivalence.com"

	// Encrypt data
	// token, err := paseto.Encrypt(symmetricKey, jsonToken, footer)
	token, err := paseto.NewV2().Encrypt(symmetricKey, jsonToken, footer)

	if err != nil {
		panic(err)
	}

	return token, nil

}
