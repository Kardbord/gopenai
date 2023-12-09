package authentication_test

import (
	"testing"

	"github.com/Kardbord/gopenai/authentication"
)

func TestAuthentication(t *testing.T) {
	err := authentication.SetAPIKey("")
	if err == nil {
		t.Fatalf("Expected an error due to empty-string API key.")
	}

	err = authentication.SetAPIKey("Your Key Here")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
