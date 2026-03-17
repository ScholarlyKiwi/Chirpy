package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	type testCase struct {
		tokenSecret string
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
		tokenString, err := MakeJWT(id, secret)
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

func TestJWTRepeating(t *testing.T) {

	id, _ := uuid.NewUUID()
	secret := "Secret"
	tokenString, _ := MakeJWT(id, secret)
	time.Sleep(2 * time.Second)
	tokenString2, _ := MakeJWT(id, secret)

	if tokenString == tokenString2 {
		t.Error("Epecting tokens to be different")
	}
}

func TestDecodeJWT(t *testing.T) {
	id, _ := uuid.NewUUID()
	secret := "Secret"
	tokenString, err := MakeJWT(id, secret)
	if err != nil {
		t.Error(err)
	}
	decodedJWT, claims, err := DecodeJWT(tokenString, secret)
	if err != nil {
		t.Error(err)
	}
	if !decodedJWT.Valid {
		t.Error("DecodeJWT returned invalid flag")
	}
	subject, err := claims.GetSubject()
	if subject != id.String() {
		t.Error("DecodeJWT ID mismatch from claims")
	}

}

/*  Test inavlid as MakeJWT no longer takes an ExpiresIn value.
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
*/
