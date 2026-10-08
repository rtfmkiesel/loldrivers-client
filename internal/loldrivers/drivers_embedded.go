//go:build embedded

package loldrivers

import (
	_ "embed"

	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
)

//go:generate curl -O https://www.loldrivers.io/api/drivers.json

var (
	//go:embed drivers.json
	embeddedDriversJson []byte
)

// In the embedded version, we bundle a drivers.json
// at build time using go:embed. This might flag
// the binary as malware. See
// https://github.com/rtfmkiesel/loldrivers-client/issues/4
//
// Most likely this is because of specific strings in drivers.json as
// such embedded resources can be read directly. (not like regular Go strings)
func getRawDrivers() ([]byte, error) {
	logger.Debug("Loading embedded drivers")
	return embeddedDriversJson, nil
}
