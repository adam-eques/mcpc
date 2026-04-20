package main

import "testing"

func TestParseArgs(t *testing.T) {
	m, err := parseArgs([]string{"expression=2+2", "limit:=5", "flag:=true"})
	if err != nil {
		t.Fatal(err)
	}
	if m["expression"] != "2+2" {
		t.Fatalf("expression=%v", m["expression"])
	}
	if m["limit"].(float64) != 5 {
		t.Fatalf("limit=%v", m["limit"])
	}
	if m["flag"] != true {
		t.Fatalf("flag=%v", m["flag"])
	}
}

func TestParseArgsErrors(t *testing.T) {
	if _, err := parseArgs([]string{"noequals"}); err == nil {
		t.Fatal("expected error for token without =")
	}
	if _, err := parseArgs([]string{"k:=not json"}); err == nil {
		t.Fatal("expected error for invalid json")
	}
}
