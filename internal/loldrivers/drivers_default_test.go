//go:build !embedded

package loldrivers

import "testing"

func TestOnlineParse(t *testing.T) {
	if err := LoadDrivers("online", ""); err != nil {
		t.Error(err)
	}
}
