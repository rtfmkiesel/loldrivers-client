package loldrivers

import "testing"

// online
func TestOnlineParse(t *testing.T) {
	jsonBytes, err := downloadNewestData()
	if err != nil {
		t.Error(err)
	}

	if err := loadJsonIntoHashmaps(jsonBytes); err != nil {
		t.Error(err)
	}
}

// internal
func TestInternalParse(t *testing.T) {
	if err := loadJsonIntoHashmaps(internalDrivers); err != nil {
		t.Error(err)
	}
}
