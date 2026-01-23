package model

// BookLevel represents one price level in the order book.
type BookLevel struct {
	Price    Price
	Quantity Quantity
}

// IsZero returns true if the level has no data.
func (l BookLevel) IsZero() bool {
	return l.Price == 0 && l.Quantity == 0
}
