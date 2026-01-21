package model

type Balance struct {
	Currency    string
	Cash        int64
	BuyingPower int64
	Raw         Attributes
}
