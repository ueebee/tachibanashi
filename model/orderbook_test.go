package model

import "testing"

func TestBookLevel_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		level BookLevel
		want  bool
	}{
		{"zero", BookLevel{}, true},
		{"has price", BookLevel{Price: 100, Quantity: 0}, false},
		{"has quantity", BookLevel{Price: 0, Quantity: 10}, false},
		{"has both", BookLevel{Price: 100, Quantity: 10}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderBook_Spread(t *testing.T) {
	tests := []struct {
		name string
		ob   OrderBook
		want Price
	}{
		{
			name: "normal spread",
			ob: OrderBook{
				Asks: [10]BookLevel{{Price: 1005, Quantity: 100}},
				Bids: [10]BookLevel{{Price: 1000, Quantity: 100}},
			},
			want: 5,
		},
		{
			name: "zero ask price",
			ob: OrderBook{
				Asks: [10]BookLevel{{Price: 0, Quantity: 100}},
				Bids: [10]BookLevel{{Price: 1000, Quantity: 100}},
			},
			want: 0,
		},
		{
			name: "zero bid price",
			ob: OrderBook{
				Asks: [10]BookLevel{{Price: 1005, Quantity: 100}},
				Bids: [10]BookLevel{{Price: 0, Quantity: 100}},
			},
			want: 0,
		},
		{
			name: "both zero",
			ob:   OrderBook{},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ob.Spread(); got != tt.want {
				t.Errorf("Spread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderBook_MidPrice(t *testing.T) {
	tests := []struct {
		name string
		ob   OrderBook
		want Price
	}{
		{
			name: "normal mid price",
			ob: OrderBook{
				Asks: [10]BookLevel{{Price: 1010, Quantity: 100}},
				Bids: [10]BookLevel{{Price: 1000, Quantity: 100}},
			},
			want: 1005,
		},
		{
			name: "odd sum",
			ob: OrderBook{
				Asks: [10]BookLevel{{Price: 1011, Quantity: 100}},
				Bids: [10]BookLevel{{Price: 1000, Quantity: 100}},
			},
			want: 1005, // (1011 + 1000) / 2 = 1005 (integer division)
		},
		{
			name: "zero ask price",
			ob: OrderBook{
				Asks: [10]BookLevel{{Price: 0, Quantity: 100}},
				Bids: [10]BookLevel{{Price: 1000, Quantity: 100}},
			},
			want: 0,
		},
		{
			name: "zero bid price",
			ob: OrderBook{
				Asks: [10]BookLevel{{Price: 1010, Quantity: 100}},
				Bids: [10]BookLevel{{Price: 0, Quantity: 100}},
			},
			want: 0,
		},
		{
			name: "both zero",
			ob:   OrderBook{},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ob.MidPrice(); got != tt.want {
				t.Errorf("MidPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
