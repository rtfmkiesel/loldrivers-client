package loldrivers

import "testing"

// online
func TestOnlineParse(t *testing.T) {
	jsonBytes, err := download()
	if err != nil {
		t.Error(err)
	}

	_, err = parse(jsonBytes)
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
