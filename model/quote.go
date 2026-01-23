package model

// Quote is a generic snapshot for a single symbol.
type Quote struct {
	Symbol string
	Fields Attributes
}

func (q Quote) Value(key string) string {
	return q.Fields.Value(key)
}

func (q Quote) Int64(key string) (int64, bool) {
	return q.Fields.Int64(key)
}

func (q Quote) Float64(key string) (float64, bool) {
	return q.Fields.Float64(key)
}

func (q Quote) Price(key string) (Price, bool) {
	value, ok := q.Fields.Int64(key)
	if !ok {
		return 0, false
	}
	return Price(value), true
}

func (q Quote) Quantity(key string) (Quantity, bool) {
	value, ok := q.Fields.Int64(key)
	if !ok {
		return 0, false
	}
	return Quantity(value), true
}

func (q Quote) LastPrice() (Price, bool) {
	return q.Price(FieldLastPrice)
}

func (q Quote) LastTime() string {
	return q.Value(FieldLastTime)
}

func (q Quote) PrevClose() (Price, bool) {
	return q.Price(FieldPrevClose)
}

func (q Quote) OpenPrice() (Price, bool) {
	return q.Price(FieldOpenPrice)
}

func (q Quote) OpenTime() string {
	return q.Value(FieldOpenTime)
}

func (q Quote) HighPrice() (Price, bool) {
	return q.Price(FieldHighPrice)
}

func (q Quote) HighTime() string {
	return q.Value(FieldHighTime)
}

func (q Quote) LowPrice() (Price, bool) {
	return q.Price(FieldLowPrice)
}

func (q Quote) LowTime() string {
	return q.Value(FieldLowTime)
}

func (q Quote) Volume() (Quantity, bool) {
	return q.Quantity(FieldVolume)
}

func (q Quote) Turnover() (int64, bool) {
	return q.Int64(FieldTurnover)
}

func (q Quote) Change() (Price, bool) {
	return q.Price(FieldChange)
}

func (q Quote) ChangeRate() (float64, bool) {
	return q.Float64(FieldChangeRate)
}

func (q Quote) VWAP() (Price, bool) {
	return q.Price(FieldVWAP)
}

func (q Quote) PrevCompare() string {
	return q.Value(FieldPrevCompare)
}

func (q Quote) BestAsk() (Price, bool) {
	return q.Price(FieldAskPrice)
}

func (q Quote) BestAskSize() (Quantity, bool) {
	return q.Quantity(FieldAskSize)
}

func (q Quote) BestAskType() string {
	return q.Value(FieldAskType)
}

func (q Quote) BestBid() (Price, bool) {
	return q.Price(FieldBidPrice)
}

func (q Quote) BestBidSize() (Quantity, bool) {
	return q.Quantity(FieldBidSize)
}

func (q Quote) BestBidType() string {
	return q.Value(FieldBidType)
}

func (q Quote) MarketAskSize() (Quantity, bool) {
	return q.Quantity(FieldMarketAskSize)
}

func (q Quote) MarketBidSize() (Quantity, bool) {
	return q.Quantity(FieldMarketBidSize)
}

// OrderBook returns 10-level order book from quote fields.
func (q Quote) OrderBook() OrderBook {
	var ob OrderBook

	askPriceFields := [10]string{
		FieldAskPrice1, FieldAskPrice2, FieldAskPrice3, FieldAskPrice4, FieldAskPrice5,
		FieldAskPrice6, FieldAskPrice7, FieldAskPrice8, FieldAskPrice9, FieldAskPrice10,
	}
	askSizeFields := [10]string{
		FieldAskSize1, FieldAskSize2, FieldAskSize3, FieldAskSize4, FieldAskSize5,
		FieldAskSize6, FieldAskSize7, FieldAskSize8, FieldAskSize9, FieldAskSize10,
	}
	bidPriceFields := [10]string{
		FieldBidPrice1, FieldBidPrice2, FieldBidPrice3, FieldBidPrice4, FieldBidPrice5,
		FieldBidPrice6, FieldBidPrice7, FieldBidPrice8, FieldBidPrice9, FieldBidPrice10,
	}
	bidSizeFields := [10]string{
		FieldBidSize1, FieldBidSize2, FieldBidSize3, FieldBidSize4, FieldBidSize5,
		FieldBidSize6, FieldBidSize7, FieldBidSize8, FieldBidSize9, FieldBidSize10,
	}

	for i := 0; i < 10; i++ {
		if price, ok := q.Price(askPriceFields[i]); ok {
			ob.Asks[i].Price = price
		}
		if qty, ok := q.Quantity(askSizeFields[i]); ok {
			ob.Asks[i].Quantity = qty
		}
		if price, ok := q.Price(bidPriceFields[i]); ok {
			ob.Bids[i].Price = price
		}
		if qty, ok := q.Quantity(bidSizeFields[i]); ok {
			ob.Bids[i].Quantity = qty
		}
	}

	return ob
}
