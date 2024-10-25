package dto

type ZigZagType int

const (
	ZigZagHigh ZigZagType = iota
	ZigZagLow
)

type ZigZagPoint struct {
	Candle Candle
	Type   ZigZagType
	Price  float64
}

func (zz ZigZagPoint) String() string {
	switch zz.Type {
	case ZigZagHigh:
		return "High"
	// case ZigZagLow:
	default:
		return "Low"
	}
}
