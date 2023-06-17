package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/o1egl/paseto"
	"github.com/xatta-trone/words-combinator/model"
)

type AuthInterface interface {
	 GenerateTokenFromEmail(user model.UserModel) (string, time.Time, error)
	//  ValidateToken(s string) string
}

type AuthService struct {

}

func NewAuthService()*AuthService{
	return &AuthService{}
}


func (a *AuthService) GenerateTokenFromEmail(user model.UserModel) (string, time.Time, error) {

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
		Jti:        user.Email,
		Subject:    strconv.Itoa(int(user.ID)),
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
		 
	}
	// Add custom claim    to the token
	userData,_ := json.Marshal(user)
	jsonToken.Set("email", user.Email)
	jsonToken.Set("user_id", strconv.Itoa(int(user.ID)))
	jsonToken.Set("user",string(userData))
	footer := "gre-sentence-equivalence.com"

	// Encrypt data
	// token, err := paseto.Encrypt(symmetricKey, jsonToken, footer)
	token, err := paseto.NewV2().Encrypt(symmetricKey, jsonToken, footer)

	if err != nil {
		panic(err)
	}

	return token,exp, nil

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
