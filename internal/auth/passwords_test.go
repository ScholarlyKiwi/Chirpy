package auth

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHashPassword(t *testing.T) {
	testCases := []string{
		"Password",
		"Secret",
		"T0pSecret!",
	}
	for _, testCase := range testCases {
		hash, err := HashPassword(testCase)
		if err != nil {
			t.Errorf("Error hashing password: %v", err)
		}
		ok, err := CheckPasswordHash(testCase, hash)
		if err != nil {
			t.Errorf("Error checking password: %v", err)
		}
		if !ok {
			t.Errorf("Password check failed: %v, %v", testCase, hash)
		}
	}
}
