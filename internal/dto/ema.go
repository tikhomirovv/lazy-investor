package dto

import "time"

type EMA struct {
	Name   string
	Dates  []time.Time
	Values []float64
}
