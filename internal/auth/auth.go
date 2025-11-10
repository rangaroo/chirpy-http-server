package auth

import (
	"net/http"
	"fmt"
	"strings"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/alexedwards/argon2id"
)

func GetBearerToken(headers http.Header) (string, error) {
	headerString := headers.Get("Authorization")

	splitted := strings.Fields(headerString)
	if len(splitted) != 2  || splitted[0] != "Bearer" {
		return "", fmt.Errorf("Wrong header format: 'Bearer <token_string>'")
	}

	return splitted[1], nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
}

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
