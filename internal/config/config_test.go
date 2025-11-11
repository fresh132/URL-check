package config_test

import (
	"testing"

	"github.com/fresh132/URL-check/internal/config"
)

func TestCheckURL_Basic(t *testing.T) {
	res := config.CheckURL("https://example.com")
	if res != "Подключение успешно выполнено" && res != "Недоступен" {
		t.Errorf("unexpected result: %s", res)
	}
}

func TestCheckURL_Empty(t *testing.T) {
	if config.CheckURL("") != "Недоступен" {
		t.Error("expected 'Недоступен' for empty URL")
	}
}
