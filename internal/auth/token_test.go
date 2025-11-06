package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-123"
	token, err := MakeJWT(userID, secret, 15 * time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}
	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error for valid token: %v", err)
	}
	if gotID != userID {
		t.Fatalf("ValidateJWT returned wrong user id: got %v want %v", gotID, userID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-expired"
	// Create a token that is already expired
	token, err := MakeJWT(userID, secret, 2 * time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}

	time.Sleep(4 * time.Second)
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("ValidateJWT did not return an error for an expired token")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	userID := uuid.New()
	secret := "correct-secret"
	wrongSecret := "wrong-secret"
	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("ValidateJWT did not return an error for a token signed with a different secret")
	}
}

func TestMalformedJWT(t *testing.T) {
	_, err := ValidateJWT("this-is-not-a-jwt", "any-secret")
	if err == nil {
		t.Fatal("ValidateJWT did not return an error for a malformed token")
	}
}
