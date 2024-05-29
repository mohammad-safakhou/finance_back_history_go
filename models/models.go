package models

import "time"

// Candle represents a single candlestick data point
type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type Trade struct {
	IsOpen     bool
	EntryTime  time.Time
	EntryPrice float64
	ExitTime   time.Time
	ExitPrice  float64
	Type       string // "long" or "short"
}
