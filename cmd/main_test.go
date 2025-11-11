package main

import (
	"os"
	"testing"
)

func TestEnvPort(t *testing.T) {
	os.Setenv("PORT", "3000")
	if got := os.Getenv("PORT"); got != "3000" {
		t.Fatalf("expected 3000, got %s", got)
	}
}
