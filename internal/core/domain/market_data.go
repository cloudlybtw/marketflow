package domain

type Message struct {
	Exchange  string          `json:"exchange"`
	Timestamp time.Time       `json:"timestamp"`
	Raw       json.RawMessage `json:"raw"`
}
