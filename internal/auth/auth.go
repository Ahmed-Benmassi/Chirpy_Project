package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess  TokenType = "chirpy-access"

)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// HashPassword -
func HashPassword(password string) (string, error) {                                  // HashPassword takes a plaintext password and returns a hashed version of it using the argon2id algorithm. It returns an error if the hashing process fails.
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil                                                           // return the hashed password if successful
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) (bool, error) {               // CheckPasswordHash compares a plaintext password with a hashed password and returns true if they match, false otherwise. It uses the argon2id package to perform the comparison.
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil                                                       // return whether the password matches the hash
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {     // MakeJWT creates a JWT token string for the given user ID, secret, and expiration duration.
	signingKey:=[]byte(tokenSecret)
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{               // RegisteredClaims is a struct that includes standard JWT claims such as Issuer, IssuedAt, ExpiresAt, and Subject.
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)                                           // return the signed token string
	
}


func ValidateJWT(tokenString,tokenSecret string)(uuid.UUID,error){          // ValidateJWT validates a JWT token string using the provided secret and returns the user ID if the token is valid.
	claimsStruct:=jwt.RegisteredClaims{}
	token,err:=jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil},
	)

	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil                                                    // return the user ID if the token is valid
}

func GetBearerToken(headers http.Header) (string, error) {// GetBearerToken extracts the Bearer token from the Authorization header of an HTTP request.
	authHeader := headers.Get("Authorization")
	if authHeader== ""{
		return "",ErrNoAuthHeaderIncluded
	}
	splitAuth :=strings.Split(authHeader," ")
	if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {     // check if the header is in the format
	    return "",errors.New("malformed authorization header")
	
	}

	return splitAuth[1], nil     // return the token part of the header
}

// MakeRefreshToken makes a random 256 bit token
// encoded in hex
func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}


// GetAPIKey -
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}