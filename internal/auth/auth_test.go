package auth

import (
	"testing"

	"net/http"
	"time"
	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
		wantToken   string
		wantErr     bool
	}{
		{
			name:        "valid token",
			headerValue: "Bearer abc.def.ghi",
			wantToken:   "abc.def.ghi",
			wantErr:     false,
		},
		{
			name:        "valid token extra spaces",
			headerValue: "Bearer    token123   ",
			wantToken:   "token123",
			wantErr:     false,
		},
		{
			name:        "missing header",
			headerValue: "",
			wantToken:   "",
			wantErr:     true,
		},
		{
			name:        "wrong scheme",
			headerValue: "Basic token",
			wantToken:   "",
			wantErr:     true,
		},
		{
			name:        "empty token after bearer",
			headerValue: "Bearer ",
			wantToken:   "",
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := http.Header{}
			if tc.headerValue != "" {
				h.Set("Authorization", tc.headerValue)
			}

			got, err := GetBearerToken(h)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil, token=%q", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantToken {
				t.Fatalf("token mismatch: got %q want %q", got, tc.wantToken)
			}
		})
	}
}

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
	token, err := MakeJWT(userID, secret, time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}

	time.Sleep(2 * time.Second)
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

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

