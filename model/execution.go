package model

type Execution struct {
	OrderID  string
	Symbol   string
	Price    Price
	Quantity Quantity
	Time     string
	Raw      Attributes
}
