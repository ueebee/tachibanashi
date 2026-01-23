package model

// Field keys for quote snapshots (sTargetColumn).
// See docs/info_code_mapping.md for the full mapping table.
const (
	FieldLastPrice = "pDPP"
	FieldLastTime  = "tDPP:T"
	FieldPrevClose = "pPRP"

	FieldOpenPrice = "pDOP"
	FieldOpenTime  = "tDOP:T"
	FieldHighPrice = "pDHP"
	FieldHighTime  = "tDHP:T"
	FieldLowPrice  = "pDLP"
	FieldLowTime   = "tDLP:T"

	FieldVolume      = "pDV"
	FieldTurnover    = "pDJ"
	FieldChange      = "pDYWP"
	FieldChangeRate  = "pDYRP"
	FieldVWAP        = "pVWAP"
	FieldPrevCompare = "pDPG"

	FieldAskPrice = "pQAP"
	FieldAskSize  = "pAV"
	FieldAskType  = "pQAS"
	FieldBidPrice = "pQBP"
	FieldBidSize  = "pBV"
	FieldBidType  = "pQBS"

	FieldMarketAskSize = "pAAV"
	FieldMarketBidSize = "pABV"

	// 10-level order book ask prices
	FieldAskPrice1  = "pGAP1"
	FieldAskPrice2  = "pGAP2"
	FieldAskPrice3  = "pGAP3"
	FieldAskPrice4  = "pGAP4"
	FieldAskPrice5  = "pGAP5"
	FieldAskPrice6  = "pGAP6"
	FieldAskPrice7  = "pGAP7"
	FieldAskPrice8  = "pGAP8"
	FieldAskPrice9  = "pGAP9"
	FieldAskPrice10 = "pGAP10"

	// 10-level order book bid prices
	FieldBidPrice1  = "pGBP1"
	FieldBidPrice2  = "pGBP2"
	FieldBidPrice3  = "pGBP3"
	FieldBidPrice4  = "pGBP4"
	FieldBidPrice5  = "pGBP5"
	FieldBidPrice6  = "pGBP6"
	FieldBidPrice7  = "pGBP7"
	FieldBidPrice8  = "pGBP8"
	FieldBidPrice9  = "pGBP9"
	FieldBidPrice10 = "pGBP10"

	// 10-level order book ask sizes
	FieldAskSize1  = "pGAV1"
	FieldAskSize2  = "pGAV2"
	FieldAskSize3  = "pGAV3"
	FieldAskSize4  = "pGAV4"
	FieldAskSize5  = "pGAV5"
	FieldAskSize6  = "pGAV6"
	FieldAskSize7  = "pGAV7"
	FieldAskSize8  = "pGAV8"
	FieldAskSize9  = "pGAV9"
	FieldAskSize10 = "pGAV10"

	// 10-level order book bid sizes
	FieldBidSize1  = "pGBV1"
	FieldBidSize2  = "pGBV2"
	FieldBidSize3  = "pGBV3"
	FieldBidSize4  = "pGBV4"
	FieldBidSize5  = "pGBV5"
	FieldBidSize6  = "pGBV6"
	FieldBidSize7  = "pGBV7"
	FieldBidSize8  = "pGBV8"
	FieldBidSize9  = "pGBV9"
	FieldBidSize10 = "pGBV10"
)
