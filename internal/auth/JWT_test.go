package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	type testCase struct {
		tokenSecret string
	}

	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("Error parsing duration %v", err)
		return
	}
	testCases := []testCase{
		{
			tokenSecret: "T0psecret!",
		},
		{
			tokenSecret: "JJW@()FF",
		},
		{
			tokenSecret: "Ikki",
		},
	}

	for idx := range testCases {
		id, _ := uuid.NewUUID()
		secret := testCases[idx].tokenSecret
		tokenString, err := MakeJWT(id, secret, duration)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		returnedId, err := ValidateJWT(tokenString, secret)
		if err != nil {
			t.Errorf("Error validating token %v", err)
			continue
		}
		if returnedId != id {
			t.Errorf("Error in JWT expected ID: %v\nResult ID: %v\n", id, returnedId)
		}
	}
}

func TestJWTExpires(t *testing.T) {
	type testCase struct {
		tokenSecret string
	}

	duration, err := time.ParseDuration("1ms")
	if err != nil {
		t.Errorf("Error parsing duration %v", err)
		return
	}
	testCases := []testCase{
		{
			tokenSecret: "T0psecret!",
		},
		{
			tokenSecret: "JJW@()FF",
		},
		{
			tokenSecret: "Ikki",
		},
	}

	for idx := range testCases {
		id, _ := uuid.NewUUID()
		secret := testCases[idx].tokenSecret
		tokenString, err := MakeJWT(id, secret, duration)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		time.Sleep(5 * time.Second)
		_, err = ValidateJWT(tokenString, secret)

		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				continue
			} else {
				t.Errorf("Error validating token %v", err)
				continue
			}
		} else {
			t.Errorf("Error - token should expire.")
		}

	}
}
