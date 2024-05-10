package dto

type TrendType int

const (
	TrendUp TrendType = iota
	TrendDown
	TrendNo // боковик, консолидация
)

type TrendChange struct {
	// Candle Candle
	Swing Swing
	Trend TrendType
}

func (tt TrendType) String() string {
	switch tt {
	case TrendUp:
		return "Up"
	case TrendDown:
		return "Down"
	default:
		return "No"
	}
}
