//go:build embedded

package loldrivers

import (
	_ "embed"
)

var (
	//go:embed drivers.json
	embeddedDriversJson []byte
)

func getEmbeddedDrivers() ([]byte, error) {
	return embeddedDriversJson, nil
}
