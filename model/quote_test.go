package model

import "testing"

func TestQuote_OrderBook_Full(t *testing.T) {
	// Test with full 10 levels of order book data
	fields := Attributes{
		// Ask prices (pGAP1-10)
		"pGAP1":  "1001",
		"pGAP2":  "1002",
		"pGAP3":  "1003",
		"pGAP4":  "1004",
		"pGAP5":  "1005",
		"pGAP6":  "1006",
		"pGAP7":  "1007",
		"pGAP8":  "1008",
		"pGAP9":  "1009",
		"pGAP10": "1010",
		// Ask quantities (pGAV1-10)
		"pGAV1":  "100",
		"pGAV2":  "200",
		"pGAV3":  "300",
		"pGAV4":  "400",
		"pGAV5":  "500",
		"pGAV6":  "600",
		"pGAV7":  "700",
		"pGAV8":  "800",
		"pGAV9":  "900",
		"pGAV10": "1000",
		// Bid prices (pGBP1-10)
		"pGBP1":  "999",
		"pGBP2":  "998",
		"pGBP3":  "997",
		"pGBP4":  "996",
		"pGBP5":  "995",
		"pGBP6":  "994",
		"pGBP7":  "993",
		"pGBP8":  "992",
		"pGBP9":  "991",
		"pGBP10": "990",
		// Bid quantities (pGBV1-10)
		"pGBV1":  "110",
		"pGBV2":  "220",
		"pGBV3":  "330",
		"pGBV4":  "440",
		"pGBV5":  "550",
		"pGBV6":  "660",
		"pGBV7":  "770",
		"pGBV8":  "880",
		"pGBV9":  "990",
		"pGBV10": "1100",
	}

	q := Quote{Symbol: "1234", Fields: fields}
	ob := q.OrderBook()

	// Verify all ask levels
	expectedAskPrices := []Price{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010}
	expectedAskQuantities := []Quantity{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
	for i := 0; i < 10; i++ {
		if ob.Asks[i].Price != expectedAskPrices[i] {
			t.Errorf("Asks[%d].Price = %d, want %d", i, ob.Asks[i].Price, expectedAskPrices[i])
		}
		if ob.Asks[i].Quantity != expectedAskQuantities[i] {
			t.Errorf("Asks[%d].Quantity = %d, want %d", i, ob.Asks[i].Quantity, expectedAskQuantities[i])
		}
	}

	// Verify all bid levels
	expectedBidPrices := []Price{999, 998, 997, 996, 995, 994, 993, 992, 991, 990}
	expectedBidQuantities := []Quantity{110, 220, 330, 440, 550, 660, 770, 880, 990, 1100}
	for i := 0; i < 10; i++ {
		if ob.Bids[i].Price != expectedBidPrices[i] {
			t.Errorf("Bids[%d].Price = %d, want %d", i, ob.Bids[i].Price, expectedBidPrices[i])
		}
		if ob.Bids[i].Quantity != expectedBidQuantities[i] {
			t.Errorf("Bids[%d].Quantity = %d, want %d", i, ob.Bids[i].Quantity, expectedBidQuantities[i])
		}
	}

	// Verify spread calculation works
	spread := ob.Spread()
	if spread != 2 {
		t.Errorf("Spread() = %d, want 2", spread)
	}
}

func TestQuote_OrderBook_Partial(t *testing.T) {
	// Test with only some levels populated
	fields := Attributes{
		// Only 3 levels of asks
		"pGAP1": "1001",
		"pGAP2": "1002",
		"pGAP3": "1003",
		"pGAV1": "100",
		"pGAV2": "200",
		"pGAV3": "300",
		// Only 2 levels of bids
		"pGBP1": "999",
		"pGBP2": "998",
		"pGBV1": "110",
		"pGBV2": "220",
	}

	q := Quote{Symbol: "1234", Fields: fields}
	ob := q.OrderBook()

	// Check populated ask levels
	if ob.Asks[0].Price != 1001 || ob.Asks[0].Quantity != 100 {
		t.Errorf("Asks[0] = {%d, %d}, want {1001, 100}", ob.Asks[0].Price, ob.Asks[0].Quantity)
	}
	if ob.Asks[1].Price != 1002 || ob.Asks[1].Quantity != 200 {
		t.Errorf("Asks[1] = {%d, %d}, want {1002, 200}", ob.Asks[1].Price, ob.Asks[1].Quantity)
	}
	if ob.Asks[2].Price != 1003 || ob.Asks[2].Quantity != 300 {
		t.Errorf("Asks[2] = {%d, %d}, want {1003, 300}", ob.Asks[2].Price, ob.Asks[2].Quantity)
	}

	// Check empty ask levels are zero
	for i := 3; i < 10; i++ {
		if !ob.Asks[i].IsZero() {
			t.Errorf("Asks[%d] should be zero, got {%d, %d}", i, ob.Asks[i].Price, ob.Asks[i].Quantity)
		}
	}

	// Check populated bid levels
	if ob.Bids[0].Price != 999 || ob.Bids[0].Quantity != 110 {
		t.Errorf("Bids[0] = {%d, %d}, want {999, 110}", ob.Bids[0].Price, ob.Bids[0].Quantity)
	}
	if ob.Bids[1].Price != 998 || ob.Bids[1].Quantity != 220 {
		t.Errorf("Bids[1] = {%d, %d}, want {998, 220}", ob.Bids[1].Price, ob.Bids[1].Quantity)
	}

	// Check empty bid levels are zero
	for i := 2; i < 10; i++ {
		if !ob.Bids[i].IsZero() {
			t.Errorf("Bids[%d] should be zero, got {%d, %d}", i, ob.Bids[i].Price, ob.Bids[i].Quantity)
		}
	}
}

func TestQuote_OrderBook_Empty(t *testing.T) {
	// Test with no order book data
	q := Quote{Symbol: "1234", Fields: Attributes{}}
	ob := q.OrderBook()

	// All levels should be zero
	for i := 0; i < 10; i++ {
		if !ob.Asks[i].IsZero() {
			t.Errorf("Asks[%d] should be zero, got {%d, %d}", i, ob.Asks[i].Price, ob.Asks[i].Quantity)
		}
		if !ob.Bids[i].IsZero() {
			t.Errorf("Bids[%d] should be zero, got {%d, %d}", i, ob.Bids[i].Price, ob.Bids[i].Quantity)
		}
	}

	// Spread should be 0 when no data
	if ob.Spread() != 0 {
		t.Errorf("Spread() = %d, want 0", ob.Spread())
	}
}

func TestQuote_OrderBook_NilFields(t *testing.T) {
	// Test with nil fields
	q := Quote{Symbol: "1234", Fields: nil}
	ob := q.OrderBook()

	// All levels should be zero
	for i := 0; i < 10; i++ {
		if !ob.Asks[i].IsZero() {
			t.Errorf("Asks[%d] should be zero, got {%d, %d}", i, ob.Asks[i].Price, ob.Asks[i].Quantity)
		}
		if !ob.Bids[i].IsZero() {
			t.Errorf("Bids[%d] should be zero, got {%d, %d}", i, ob.Bids[i].Price, ob.Bids[i].Quantity)
		}
	}
}
