package config

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)

	if err != nil {
		return "Недоступен"
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode >= 400 {

		req2, err2 := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err2 != nil {
			return "Недоступен"
		}

		res2, err2 := http.DefaultClient.Do(req2)

		if err2 != nil || res2.StatusCode >= 400 {
			return "Недоступен"
		}

		return "Подключение успешно выполнено"
	}

	return "Подключение успешно выполнено"
}
