//go:build !embedded

package loldrivers

import "fmt"

// In the normal version, we do not bundle the JSON.
// https://github.com/rtfmkiesel/loldrivers-client/issues/4
func getEmbeddedDrivers() ([]byte, error) {
	return nil, fmt.Errorf("This version does not support '-mode embedded'")
}
