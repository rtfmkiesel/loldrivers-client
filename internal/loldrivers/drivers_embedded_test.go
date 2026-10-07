//go:build embedded

package loldrivers

import "testing"

func TestEmbeddedlParse(t *testing.T) {
	if err := LoadDrivers("embedded", ""); err != nil {
		t.Error(err)
	}
}
