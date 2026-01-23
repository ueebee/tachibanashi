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

// OrderBook represents a 10-level order book.
type OrderBook struct {
	Asks [10]BookLevel // sell side (ascending price order)
	Bids [10]BookLevel // buy side (descending price order)
}

// Spread returns the difference between the best ask and best bid price.
// Returns 0 if either price is zero.
func (ob OrderBook) Spread() Price {
	if ob.Asks[0].Price == 0 || ob.Bids[0].Price == 0 {
		return 0
	}
	return ob.Asks[0].Price - ob.Bids[0].Price
}

// MidPrice returns the midpoint between the best ask and best bid price.
// Returns 0 if either price is zero.
func (ob OrderBook) MidPrice() Price {
	if ob.Asks[0].Price == 0 || ob.Bids[0].Price == 0 {
		return 0
	}
	return (ob.Asks[0].Price + ob.Bids[0].Price) / 2
}
