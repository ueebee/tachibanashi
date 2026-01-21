package model

type Position struct {
	Symbol    string
	Quantity  Quantity
	AvgPrice  Price
	UnrealPnL int64
	Raw       Attributes
}
