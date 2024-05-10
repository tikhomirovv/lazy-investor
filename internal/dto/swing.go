package dto

type SwingType int

const (
	SwingHigh SwingType = iota
	SwingLow
)

type Swing struct {
	Period int
	Type   SwingType
	Candle Candle
}

func (s Swing) GetValue() float64 {
	if s.Type == SwingHigh {
		return s.Candle.High
	}
	return s.Candle.Low
}
