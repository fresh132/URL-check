package config_test

import (
	"testing"

	"github.com/fresh132/URL-check/internal/config"
)

func TestCheckURL_Basic(t *testing.T) {
	res := config.CheckURL("https://example.com")
	if res != "Available" && res != "Not available" {
		t.Errorf("unexpected result: %s", res)
	}
}

func TestCheckURL_Empty(t *testing.T) {
	if config.CheckURL("") != "Not available" {
		t.Error("expected 'Not available' for empty URL")
	}
}
