//go:build embedded

package loldrivers

import (
	_ "embed"
)

//go:generate curl -O https://www.loldrivers.io/api/drivers.json

var (
	//go:embed drivers.json
	embeddedDriversJson []byte
)

func getEmbeddedDrivers() ([]byte, error) {
	return embeddedDriversJson, nil
}
