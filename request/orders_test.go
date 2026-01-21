package request

import (
	"encoding/json"
	"testing"

	"github.com/ueebee/tachibanashi/model"
)

func TestOrderRequestMarshalJSON(t *testing.T) {
	req := orderRequest{
		CommonParams: model.CommonParams{
			PNo:      "1",
			PSDDate:  "2020.01.02-03:04:05.000",
			JsonOfmt: "5",
		},
		CLMID: clmKabuNewOrder,
		Params: OrderParams{
			"sIssueCode":  "6501",
			"sOrderSuryou": "100",
			"aCLMKabuHensaiData": []map[string]string{
				{"sTategyokuNumber": "1", "sOrderSuryou": "100"},
			},
		},
	}

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if got["sCLMID"] != clmKabuNewOrder {
		t.Fatalf("sCLMID mismatch: %v", got["sCLMID"])
	}
	if got["p_no"] != "1" {
		t.Fatalf("p_no mismatch: %v", got["p_no"])
	}
	if got["p_sd_date"] != "2020.01.02-03:04:05.000" {
		t.Fatalf("p_sd_date mismatch: %v", got["p_sd_date"])
	}
	if got["sJsonOfmt"] != "5" {
		t.Fatalf("sJsonOfmt mismatch: %v", got["sJsonOfmt"])
	}
	if got["sIssueCode"] != "6501" {
		t.Fatalf("sIssueCode mismatch: %v", got["sIssueCode"])
	}
	if got["sOrderSuryou"] != "100" {
		t.Fatalf("sOrderSuryou mismatch: %v", got["sOrderSuryou"])
	}
	list, ok := got["aCLMKabuHensaiData"].([]any)
	if !ok || len(list) != 1 {
		t.Fatalf("aCLMKabuHensaiData mismatch: %v", got["aCLMKabuHensaiData"])
	}
}

func TestOrderListResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"p_no":"1",
		"p_sd_date":"2020.01.02-03:04:05.000",
		"sCLMID":"CLMOrderList",
		"sResultCode":"0",
		"sResultText":"",
		"sWarningCode":"0",
		"sWarningText":"",
		"sIssueCode":"8411",
		"sOrderSyoukaiStatus":"",
		"sSikkouDay":"",
		"aOrderList":[
			{
				"sOrderOrderNumber":"18000002",
				"sOrderIssueCode":"8411",
				"sOrderBaibaiKubun":"3",
				"sOrderOrderSuryou":"100",
				"sOrderOrderPrice":"2300.0000",
				"sOrderStatus":"FILLED"
			}
		]
	}`)

	var resp OrderListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.IssueCode != "8411" {
		t.Fatalf("issue code mismatch: %s", resp.IssueCode)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("entries length mismatch: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.OrderID != "18000002" {
		t.Fatalf("order id mismatch: %s", entry.OrderID)
	}
	if entry.Symbol != "8411" {
		t.Fatalf("symbol mismatch: %s", entry.Symbol)
	}
	order := entry.Order()
	if order.Quantity != 100 {
		t.Fatalf("order quantity mismatch: %d", order.Quantity)
	}
	if order.Price != 2300 {
		t.Fatalf("order price mismatch: %d", order.Price)
	}
	if order.Status != "FILLED" {
		t.Fatalf("order status mismatch: %s", order.Status)
	}
}

func TestOrderListDetailResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"p_no":"1",
		"p_sd_date":"2020.01.02-03:04:05.000",
		"sCLMID":"CLMOrderListDetail",
		"sResultCode":"0",
		"sResultText":"",
		"sWarningCode":"0",
		"sWarningText":"",
		"sOrderNumber":"18000002",
		"sEigyouDay":"20231018",
		"sIssueCode":"8411",
		"sOrderStatus":"FILLED"
	}`)

	var resp OrderListDetailResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.OrderNumber != "18000002" {
		t.Fatalf("order number mismatch: %s", resp.OrderNumber)
	}
	if resp.EigyouDay != "20231018" {
		t.Fatalf("eigyou day mismatch: %s", resp.EigyouDay)
	}
	if resp.IssueCode != "8411" {
		t.Fatalf("issue code mismatch: %s", resp.IssueCode)
	}
	if resp.Fields.Value("sOrderStatus") != "FILLED" {
		t.Fatalf("order status mismatch: %s", resp.Fields.Value("sOrderStatus"))
	}
}
