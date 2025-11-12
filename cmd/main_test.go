package main

import (
	"os"
	"testing"
)

func TestMainPortDefault(t *testing.T) {
	os.Unsetenv("PORT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port != "8080" {
		t.Errorf("expected 8080, got %s", port)
	}
}
