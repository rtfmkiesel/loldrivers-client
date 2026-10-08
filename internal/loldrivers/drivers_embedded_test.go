//go:build embedded

package loldrivers

import "testing"

func TestEmbeddedParse(t *testing.T) {
	if err := LoadDrivers(); err != nil {
		t.Error(err)
	}
}
