package loldrivers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
)

var (
	md5Sums  = make(map[string]*Driver)
	sha1Sums = make(map[string]*Driver)
	sha2Sums = make(map[string]*Driver)
)

const driverJsonPath = "./drivers.json"

func LoadDrivers() error {
	jsonBytes, err := func() ([]byte, error) {
		if _, err := os.Stat(driverJsonPath); err == nil {
			// ./drivers.json exists, load from disk
			logger.Debug("Loading drivers from %s", driverJsonPath)
			data, err := os.ReadFile(driverJsonPath)
			if err != nil {
				return nil, err
			}

			return data, nil
		}

		// drivers.json does not exist, fetch the raw bytes
		// either from the web or embedded
		data, err := getRawDrivers()
		if err != nil {
			return nil, err
		}

		// Save as ./drivers.json
		fh, err := os.Create(driverJsonPath)
		if err != nil {
			return nil, err
		}
		defer fh.Close()

		if _, err := fh.Write(data); err != nil {
			return nil, err
		}
		logger.Debug("Saved drivers to %s", driverJsonPath)

		return data, nil
	}()
	if err != nil {
		return fmt.Errorf("load drivers: %w", err)
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
