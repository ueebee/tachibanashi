package event

import (
	"encoding/base64"
	"testing"
)

func TestDecodeEventNSFields(t *testing.T) {
	title := "Market headline"
	body := "Market body"
	titleEncoded := base64.StdEncoding.EncodeToString([]byte(title))
	bodyEncoded := base64.StdEncoding.EncodeToString([]byte(body))

	raw := "p_no\x027\x01p_date\x022020.08.26-12:59:13.598\x01p_cmd\x02NS\x01p_PV\x02QNSD\x01p_ENO\x02166\x01p_ALT\x020\x01" +
		"p_ID\x0220200826125300_MIO1708\x01p_DT\x0220200826\x01p_TM\x02125300\x01p_CGN\x022\x01p_CGL\x02100\x03110\x01" +
		"p_GRN\x021\x01p_GRL\x023009\x01p_ISN\x022\x01p_ISL\x024519\x034568\x01p_SKF\x020\x01p_UPD\x02\x01" +
		"p_HDL\x02" + titleEncoded + "\x01p_TX\x02" + bodyEncoded

	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	ns, ok := event.(NS)
	if !ok {
		t.Fatalf("event type mismatch")
	}
	if ns.Provider != "QNSD" {
		t.Fatalf("p_PV = %s", ns.Provider)
	}
	if ns.EventNo != "166" {
		t.Fatalf("p_ENO = %s", ns.EventNo)
	}
	if ns.NewsID != "20200826125300_MIO1708" {
		t.Fatalf("p_ID = %s", ns.NewsID)
	}
	if ns.NewsDate != "20200826" {
		t.Fatalf("p_DT = %s", ns.NewsDate)
	}
	if ns.NewsTime != "125300" {
		t.Fatalf("p_TM = %s", ns.NewsTime)
	}
	if ns.CategoryCount != 2 {
		t.Fatalf("p_CGN = %d", ns.CategoryCount)
	}
	if len(ns.Categories) != 2 || ns.Categories[0] != "100" || ns.Categories[1] != "110" {
		t.Fatalf("p_CGL = %#v", ns.Categories)
	}
	if ns.GenreCount != 1 {
		t.Fatalf("p_GRN = %d", ns.GenreCount)
	}
	if len(ns.Genres) != 1 || ns.Genres[0] != "3009" {
		t.Fatalf("p_GRL = %#v", ns.Genres)
	}
	if ns.IssueCount != 2 {
		t.Fatalf("p_ISN = %d", ns.IssueCount)
	}
	if len(ns.Issues) != 2 || ns.Issues[0] != "4519" || ns.Issues[1] != "4568" {
		t.Fatalf("p_ISL = %#v", ns.Issues)
	}
	if ns.Headline != title {
		t.Fatalf("p_HDL = %s", ns.Headline)
	}
	if ns.Body != body {
		t.Fatalf("p_TX = %s", ns.Body)
	}
}
