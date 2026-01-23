package master

import (
	"testing"

	"github.com/ueebee/tachibanashi/model"
)

func TestMemoryStoreUpsertRules(t *testing.T) {
	store := NewMemoryStore()
	typ := MasterIssueMstKabu

	fields := model.Attributes{"sIssueCode": "6501", "sIssueName": "Alpha"}
	meta := UpdateMeta{Serial: 2, UpdatedAt: "20240101"}
	if !store.Upsert(typ, "6501", fields, meta) {
		t.Fatal("expected initial upsert to apply")
	}

	if store.Upsert(typ, "6501", model.Attributes{"sIssueName": "Old"}, UpdateMeta{Serial: 1}) {
		t.Fatal("expected older serial to be ignored")
	}

	if store.Upsert(typ, "6501", model.Attributes{"sIssueName": "New"}, UpdateMeta{Serial: 2, UpdatedAt: "20240102"}) == false {
		t.Fatal("expected same serial with newer timestamp to apply")
	}
	got, ok := store.Get(typ, "6501")
	if !ok || got.Fields["sIssueName"] != "New" {
		t.Fatalf("unexpected record after update: %#v", got)
	}

	if store.Upsert(typ, "6501", model.Attributes{"sIssueName": "TooOld"}, UpdateMeta{Serial: 2, UpdatedAt: "20230101"}) {
		t.Fatal("expected older timestamp to be ignored")
	}

	if !store.Upsert(typ, "6501", nil, UpdateMeta{Deleted: true}) {
		t.Fatal("expected delete to apply")
	}
	if _, ok := store.Get(typ, "6501"); ok {
		t.Fatal("expected record to be deleted")
	}
}

func TestMemoryStoreIndex(t *testing.T) {
	store := NewMemoryStore()
	typ := MasterIssueMstKabu
	store.RegisterIndex(typ, IndexSpec{Name: "issue_market", Fields: []string{"sIssueCode", "sSizyouC"}})

	fields := model.Attributes{"sIssueCode": "6501", "sSizyouC": "00"}
	store.Upsert(typ, "6501:00", fields, UpdateMeta{Serial: 1})

	value := JoinIndex("6501", "00")
	list := store.FindByIndex(typ, "issue_market", value)
	if len(list) != 1 || list[0].Key != "6501:00" {
		t.Fatalf("unexpected index result: %#v", list)
	}

	store.Delete(typ, "6501:00")
	list = store.FindByIndex(typ, "issue_market", value)
	if len(list) != 0 {
		t.Fatalf("expected index to be empty after delete: %#v", list)
	}
}
