package internal

import (
	"sync"
	"time"
)

type URLstatus struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

type TimeURL struct {
	Link []URLstatus `json:"link"`
	Time time.Time   `json:"time"`
}

var (
	Data = make(map[int]TimeURL)
	Mutx sync.Mutex
	ID   = 1
	File = "data.json"
)
