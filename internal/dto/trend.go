package dto

type TrendType string

const (
	Uptrend   TrendType = "Uptrend"
	Downtrend TrendType = "Downtrend"
	NoTrend   TrendType = "NoTrend"
)

type TrendChange struct {
	Candle   Candle
	NewTrend TrendType
}
