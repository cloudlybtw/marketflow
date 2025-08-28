package domain

import "time"

// type PriceUpdate struct {
// 	Symbol    string  `json:"symbol"`
// 	Price     float64 `json:"price"`
// 	Timestamp int64   `json:"timestamp"`
// }

type CoinStats struct {
	MinPrice  float64   `json:"min_price"`
	MaxPrice  float64   `json:"max_price"`
	AvgPrice  float64   `json:"avg_price"`
	LastPrice float64   `json:"last_price"`
	Count     int       `json:"count"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

type PriceUpdate struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Exchange  string  `json:"exchange"`
	Timestamp int64   `json:"timestamp"`
}
