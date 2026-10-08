//go:build !embedded

package loldrivers

import (
	"io"
	"net/http"

	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
)

// In the default version we download drivers.json at runtime.
func getRawDrivers() ([]byte, error) {
	logger.Debug("Downloading driver data")

	c := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.loldrivers.io/api/drivers.json", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "LOLDrivers-client")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logger.Debug("Downloaded %d bytes", len(jsonBytes))

	return jsonBytes, nil
}
