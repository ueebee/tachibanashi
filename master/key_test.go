package master

import (
	"testing"

	"github.com/ueebee/tachibanashi/model"
)

func TestMasterKey(t *testing.T) {
	tests := []struct {
		name   string
		typ    MasterType
		fields model.Attributes
		want   string
		ok     bool
	}{
		{
			name: "system_status",
			typ:  MasterSystemStatus,
			fields: model.Attributes{
				"sSystemStatusKey": "001",
			},
			want: "001",
			ok:   true,
		},
		{
			name: "yobine",
			typ:  MasterYobine,
			fields: model.Attributes{
				"sYobineTaniNumber": "101",
				"sTekiyouDay":       "20240101",
			},
			want: JoinIndex("101", "20240101"),
			ok:   true,
		},
		{
			name: "unyou_status",
			typ:  MasterUnyouStatus,
			fields: model.Attributes{
				"sSystemKouzaKubun": "102",
				"sUnyouCategory":    "01",
				"sUnyouUnit":        "0101",
				"sEigyouDayC":       "0",
				"sUnyouStatus":      "001",
				"sTaisyouGyoumu":    "04",
			},
			want: JoinIndex("102", "01", "0101", "0", "001", "04"),
			ok:   true,
		},
		{
			name: "unyou_status_kabu",
			typ:  MasterUnyouStatusKabu,
			fields: model.Attributes{
				"sSystemKouzaKubun": "102",
				"sZyouzyouSizyou":   "00",
				"sUnyouCategory":    "01",
				"sUnyouUnit":        "0101",
				"sEigyouDayC":       "0",
				"sUnyouStatus":      "001",
			},
			want: JoinIndex("102", "00", "01", "0101", "0"),
			ok:   true,
		},
		{
			name: "unyou_status_hasei",
			typ:  MasterUnyouStatusHasei,
			fields: model.Attributes{
				"sSystemKouzaKubun": "102",
				"sZyouzyouSizyou":   "01",
				"sGensisanCode":     "101",
				"sSyouhinType":      "03",
				"sUnyouCategory":    "02",
				"sUnyouUnit":        "0201",
				"sEigyouDayC":       "0",
			},
			want: JoinIndex("102", "01", "101", "03", "02", "0201", "0"),
			ok:   true,
		},
		{
			name: "issue_mst_kabu",
			typ:  MasterIssueMstKabu,
			fields: model.Attributes{
				"sIssueCode": "6501",
			},
			want: "6501",
			ok:   true,
		},
		{
			name: "issue_sizyou_mst_kabu",
			typ:  MasterIssueSizyouMstKabu,
			fields: model.Attributes{
				"sIssueCode":      "6501",
				"sZyouzyouSizyou": "00",
			},
			want: JoinIndex("6501", "00"),
			ok:   true,
		},
		{
			name: "issue_sizyou_kisei_kabu",
			typ:  MasterIssueSizyouKiseiKabu,
			fields: model.Attributes{
				"sSystemKouzaKubun": "102",
				"sIssueCode":        "6501",
				"sZyouzyouSizyou":   "00",
			},
			want: JoinIndex("102", "6501", "00"),
			ok:   true,
		},
		{
			name: "issue_sizyou_kisei_hasei",
			typ:  MasterIssueSizyouKiseiHasei,
			fields: model.Attributes{
				"sSystemKouzaKubun": "102",
				"sIssueCode":        "160060018",
				"sZyouzyouSizyou":   "01",
			},
			want: JoinIndex("102", "160060018", "01"),
			ok:   true,
		},
		{
			name: "daiyou_kakeme",
			typ:  MasterDaiyouKakeme,
			fields: model.Attributes{
				"sSystemKouzaKubun": "102",
				"sIssueCode":        "1352",
				"sTekiyouDay":       "20220422",
			},
			want: JoinIndex("102", "1352", "20220422"),
			ok:   true,
		},
		{
			name: "hosyoukin_mst",
			typ:  MasterHosyoukinMst,
			fields: model.Attributes{
				"sSystemKouzaKubun": "102",
				"sIssueCode":        "1356",
				"sZyouzyouSizyou":   "00",
				"sHenkouDay":        "20230110",
			},
			want: JoinIndex("102", "1356", "00", "20230110"),
			ok:   true,
		},
		{
			name: "order_err_reason",
			typ:  MasterOrderErrReason,
			fields: model.Attributes{
				"sErrReasonCode": "-110007",
			},
			want: "-110007",
			ok:   true,
		},
		{
			name: "issue_mst_other_with_market",
			typ:  MasterIssueMstOther,
			fields: model.Attributes{
				"sIssueCode":      "TOPIX",
				"sZyouzyouSizyou": "00",
			},
			want: JoinIndex("TOPIX", "00"),
			ok:   true,
		},
		{
			name: "issue_mst_other_no_market",
			typ:  MasterIssueMstOther,
			fields: model.Attributes{
				"sIssueCode": "TOPIX",
			},
			want: "TOPIX",
			ok:   true,
		},
		{
			name:   "missing_key",
			typ:    MasterIssueMstKabu,
			fields: model.Attributes{},
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := MasterKey(tt.typ, tt.fields)
			if ok != tt.ok {
				t.Fatalf("ok mismatch: got %v want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("value mismatch: got %q want %q", got, tt.want)
			}
		})
	}
}
