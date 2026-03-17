package auth

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestBearerToken(t *testing.T) {
	type testCase struct {
		tokenSecret string
	}

	testCases := []testCase{
		{
			tokenSecret: "Test1",
		},
		{
			tokenSecret: "Test2",
		},
	}

	for _, testCase := range testCases {
		url := "http://test"
		req, err := http.NewRequest("GET", url, bytes.NewBuffer((nil)))
		if err != nil {
			t.Errorf("TestBearerToken: %v", err)
			return
		}
		bearer := "Bearer " + testCase.tokenSecret
		req.Header.Set("Authorization", bearer)
		req.Header.Add("Accept", "application/json")

		token, err := GetBearerToken(req.Header)
		if err != nil {
			t.Errorf("TestBearerToken: %v", err)
		}
		if token != testCase.tokenSecret {
			t.Errorf("TestBearerToken: expected: %v\nactual: %v\n", testCase.tokenSecret, token)
		}

	}
}

func TestBearerTokenMissingHeader(t *testing.T) {

	url := "http://test"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer((nil)))
	if err != nil {
		t.Errorf("TestBearerTokenMissingHeader: %v", err)
		return
	}
	req.Header.Add("Accept", "application/json")

	_, err = GetBearerToken(req.Header)
	if err != nil {
		if strings.Contains(err.Error(), "missing authorization") {
			return
		} else {
			t.Errorf("TestBearerTokenMissingHeader: %v", err)
		}
	} else {
		t.Error("TestBearerTokenMissingHeader: should return missing authorisation error.")
	}
}

func TestBearerTokenMissingToken(t *testing.T) {

	url := "http://test"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer((nil)))
	if err != nil {
		t.Errorf("TestBearerTokenMissingHeader: %v", err)
		return
	}
	req.Header.Add("Accept", "application/json")
	bearer := "Bearer "
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	_, err = GetBearerToken(req.Header)
	if err != nil {
		if strings.Contains(err.Error(), "no bearer token supplied") {
			return
		} else {
			t.Errorf("TestBearerTokenMissingHeader: %v", err)
		}
	} else {
		t.Error("TestBearerTokenMissingHeader: should return no bearer token supplied error.")
	}
}
