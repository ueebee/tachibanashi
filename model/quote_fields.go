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
)
