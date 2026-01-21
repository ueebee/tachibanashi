package model

type CommonResponse struct {
	PNo        string `json:"p_no"`
	PSDDate    string `json:"p_sd_date"`
	PRVDate    string `json:"p_rv_date"`
	PErrNo     string `json:"p_errno"`
	PErr       string `json:"p_err"`
	CLMID      string `json:"sCLMID"`
	ResultCode string `json:"sResultCode,omitempty"`
	ResultText string `json:"sResultText,omitempty"`
}
