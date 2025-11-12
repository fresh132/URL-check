package config

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fresh132/URL-check/internal"
)

func Load() {
	if _, err := os.Stat(internal.File); os.IsNotExist(err) {
		return
	}

	bytes, err := os.ReadFile(internal.File)

	if err != nil {
		log.Println("Error loading data:", err)
		return
	}

	internal.Mutx.Lock()
	defer internal.Mutx.Unlock()

	if err := json.Unmarshal(bytes, &internal.Data); err != nil {
		log.Println("Data parsing error:", err)
		return
	}

	for id := range internal.Data {
		if id >= internal.ID {
			internal.ID = id + 1
		}
	}

	log.Println("Data loaded from file")
}

func Save() {
	internal.Mutx.Lock()

	bytes, err := json.Marshal(internal.Data)

	internal.Mutx.Unlock()

	if err != nil {
		log.Println("Data error:", err)
		return
	}

	tmp := internal.File + ".tmp"

	if err := os.WriteFile(tmp, bytes, 0644); err != nil {
		log.Println("Data save error:", err)
		return
	}

	if err := os.Rename(tmp, internal.File); err != nil {
		log.Println("Data move error:", err)
		return
	}

	log.Println("The data has been saved to a file")
}

func CheckURL(url string) string {
	if url == "" {
		return "Not available"
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "Not available"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)

	if err != nil {
		return "Not available"
	}

	res, err := client.Do(req)

	if err == nil && res.StatusCode < 400 {
		res.Body.Close()
		return "Available"
	}

	if res != nil {
		res.Body.Close()
	}

	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return "Not available"
	}

	res, err = client.Do(req)

	if err == nil && res.StatusCode < 400 {
		res.Body.Close()
		return "Available"
	}

	if res != nil {
		res.Body.Close()
	}

	return "Not available"
}
