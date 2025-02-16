package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestHashPassword(t *testing.T) {
	cases := []struct {
		password string
		expected error
	}{
		{password: "password1", expected: nil},
		{password: "11111", expected: nil},
	}

	for _, c := range cases {
		_, hashedSuccess := HashPassword(c.password)
		if c.expected != hashedSuccess {
			t.Errorf("for password: %s expected result: %v, but got: %v", c.password, c.expected, hashedSuccess)
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		inputPassword    string
		expectedPassword string
	}{
		{"password1", "$2a$10$jZ6nYOOA4XAi.iNsxlZICOGx6p4IUaaBdX3fK5JHufUhcq4D.lDlO"},
		{"11111", "$2a$10$ZB6EnDh.q4kmsRf.rydc0e0wWBNJlwiyO7.PR0LZq8LHwR/WTBWJi"},
	}

	for _, c := range cases {
		err := CheckPasswordHash(c.inputPassword, c.expectedPassword)
		if err != nil {
			t.Errorf("hash for password %s expected to be valid", c.inputPassword)
		}
	}
}

func TestValidateJWT(t *testing.T) {
	secretKey := "mysecretkey"
	userID := uuid.New()

	t.Run("ValidJWT", func(t *testing.T) {
		signedToken, err := MakeJWT(userID, secretKey, time.Hour)
		if err != nil {
			t.Fatalf("expected no error when making JWT, got %v", err)
		}

		idJWT, err := ValidateJWT(signedToken, secretKey)
		if err != nil {
			t.Fatalf("expected no error when validating JWT, got %v", err)
		}

		if userID != idJWT {
			t.Fatalf("expected userID %v, got %v", userID, idJWT)
		}
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		signedToken, err := MakeJWT(userID, secretKey, time.Duration(-1))
		if err != nil {
			t.Fatalf("expected no error when making JWT, got %v", err)
		}

		_, err = ValidateJWT(signedToken, secretKey)
		if err == nil {
			t.Fatalf("expected error for expired JWT, got none")
		}
	})

	t.Run("WrongSecretKey", func(t *testing.T) {
		signedToken, err := MakeJWT(userID, secretKey, time.Hour)
		if err != nil {
			t.Fatalf("expected no error when making JWT, got %v", err)
		}

		_, err = ValidateJWT(signedToken, "wrongsecretkey")
		if err == nil {
			t.Fatal("expected error with invalid signature, got none")
		}
	})
}

func TestGetBearerToken(t *testing.T) {
	t.Run("BearerTokenProvided", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", fmt.Sprintf("Bearer %s", testToken))

		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("expected no error when getting bearer token, got %v", err)
		}

		if token != testToken {
			t.Fatalf("tokens were expected to be equal")
		}
	})

	t.Run("BearerTokenNotProvided", func(t *testing.T) {
		headers := http.Header{}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Fatal("expected error getting bearer token, got none")
		}
	})

}
