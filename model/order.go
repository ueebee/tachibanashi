package model

type Order struct {
	ID       string
	Symbol   string
	Side     string
	Price    Price
	Quantity Quantity
	Status   string
	Raw      Attributes
}
