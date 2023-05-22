package loldrivers

import "testing"

// online
func TestOnlineParse(t *testing.T) {
	_, err := download()
	if err != nil {
		t.Error(err)
	}
}

// internal
func TestInternalParse(t *testing.T) {
	_, err := parse(internalDrivers)
	if err != nil {
		t.Error(err)
	}
}
