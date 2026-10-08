//go:build !embedded

package loldrivers

import "testing"

func TestOnlineParse(t *testing.T) {
	if err := LoadDrivers(); err != nil {
		t.Error(err)
	}
}
