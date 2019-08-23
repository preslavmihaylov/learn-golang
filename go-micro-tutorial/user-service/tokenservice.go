package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
)

var (

	// Define a secure authHash string used
	// as a salt when hashing our tokens.
	// Please make your own way more secure than this,
	// use a randomly generated md5 hash or something.
	authHash = []byte("mySuperSecretKeyLol")
)

// CustomClaims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type CustomClaims struct {
	User *proto.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *proto.User) (string, error)
}

type TokenService struct {
	repo *userRepository
}

// Decode a token string into a token object
func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return authHash, nil
		})
	if err != nil {
		return nil, fmt.Errorf("couldn't parse token: %s", err)
	}

	// Validate the token and return the custom claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// Encode a claim into a JWT
func (srv *TokenService) Encode(user *proto.User) (string, error) {
	expireToken := time.Now().Add(72 * time.Hour).Unix()

	// Create the Claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "shippy.user.service",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(authHash)
}
