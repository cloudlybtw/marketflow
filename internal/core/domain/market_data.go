package domain

import (
	"encoding/json"
	"time"
)

type Message struct {
	Exchange  string          `json:"exchange"`
	Timestamp time.Time       `json:"timestamp"`
	Raw       json.RawMessage `json:"raw"`
}

type ExchangeData struct {
	Coins map[string]*CoinStats `json:"coins"`
}

type AggregatedData struct {
	WindowStart time.Time                `json:"window_start"`
	WindowEnd   time.Time                `json:"window_end"`
	Exchanges   map[string]*ExchangeData `json:"exchanges"`
}
