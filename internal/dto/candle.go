package dto

import "time"

type Candle struct {
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     int64
	Time       time.Time
	IsComplete bool
}
