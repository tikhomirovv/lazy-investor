package dto

type Isin string

type Instrument struct {
	Name string
	Isin Isin
	Uid  string
}
