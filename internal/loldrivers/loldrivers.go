package loldrivers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
)

//go:generate curl -O https://www.loldrivers.io/api/drivers.json

var (
	//go:embed drivers.json
	embeddedDriversJson []byte

	md5Sums  = make(map[string]*Driver)
	sha1Sums = make(map[string]*Driver)
	sha2Sums = make(map[string]*Driver)
)

func LoadDrivers(mode string, path string) error {
	jsonBytes, err := func(mode, path string) ([]byte, error) {
		logger.Debug("Loading drivers (mode=%s)", mode)
		switch mode {
		case "online":
			return fetchDrivers()
		case "file":
			return os.ReadFile(path)
		case "embedded":
			return embeddedDriversJson, nil
		default:
			return nil, fmt.Errorf("invalid mode")
		}
	}(mode, path)
	if err != nil {
		return fmt.Errorf("load drivers: mode=%s: %w", mode, err)
	}

	drivers := []*Driver{}
	if err := json.Unmarshal(jsonBytes, &drivers); err != nil {
		return fmt.Errorf("load drivers: unmarshal json: %w", err)
	}
	logger.Debug("Loaded a total of %d drivers", len(drivers))

	for _, driver := range drivers {
		for _, sample := range driver.KnownVulnerableSamples {
			if sample.MD5 != "" {
				md5Sums[sample.MD5] = driver
			}
			if sample.SHA1 != "" {
				sha1Sums[sample.SHA1] = driver
			}
			if sample.SHA256 != "" {
				sha2Sums[sample.SHA256] = driver
			}
		}
	}
	logger.Debug("Prepared checksums (md5=%d,sha1=%d,sha2=%d)", len(md5Sums), len(sha1Sums), len(sha2Sums))

	return nil
}

func fetchDrivers() ([]byte, error) {
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

func FindDriverByHash(hash string) (bool, *Driver) {
	var d *Driver

	switch len(hash) {
	case 32:
		d = md5Sums[hash]
	case 40:
		d = sha1Sums[hash]
	case 64:
		d = sha2Sums[hash]
	default:
		return false, nil
	}

	return d != nil, d
}
