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
