package loldrivers

import "testing"

func TestOnlineParse(t *testing.T) {
	if err := LoadDrivers("online", ""); err != nil {
		t.Error(err)
	}
}

func TestEmbeddedlParse(t *testing.T) {
	if err := LoadDrivers("embedded", ""); err != nil {
		t.Error(err)
	}
}
