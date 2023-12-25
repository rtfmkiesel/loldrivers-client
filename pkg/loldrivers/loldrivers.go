package loldrivers

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
)

var (
	LoadedDrivers []Driver
	LoadedHashes  DriverHashes
)

// LoadDrivers() will load the drivers based on the selected mode
//
// mode = online, local, internal
//
// filePath = path to JSON file if mode == ModeLocal
func LoadDrivers(mode string, filePath string) (err error) {
	switch mode {
	case "online":
		jsonBytes, err := download()
		if err != nil {
			logger.Error(err)
			logger.Log("[*] Got an error while downloading data. Falling back to internal data set")
			LoadedDrivers, err = parse(internalDrivers)
			if err != nil {
				return err
			}
		} else {
			LoadedDrivers, err = parse(jsonBytes)
			if err != nil {
				logger.Error(err)
				logger.Log("[*] Got an error while parsing data. Falling back to internal data set")
				LoadedDrivers, err = parse(internalDrivers)
				if err != nil {
					return err
				}
			}
		}

	case "local":
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("could not open file '%s'", filePath)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("could not read file '%s'", filePath)
		}

		LoadedDrivers, err = parse(content)
		if err != nil {
			logger.Error(err)
			logger.Logf("[*] Got an error while parsing '%s'. Falling back to internal data set", filePath)
			LoadedDrivers, err = parse(internalDrivers)
			if err != nil {
				return err
			}
		}

	case "internal":
		LoadedDrivers, err = parse(internalDrivers)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid mode '%s'", mode)
	}

	logger.Logf("[+] Loaded %d drivers", len(LoadedDrivers))

	for _, driver := range LoadedDrivers {
		for _, knownVulnSample := range driver.KnownVulnerableSamples {
			// Append MD5 if exist
			if knownVulnSample.MD5 != "" && knownVulnSample.MD5 != "-" {
				LoadedHashes.MD5Sums = append(LoadedHashes.MD5Sums, knownVulnSample.MD5)
			}
			// Append SHA1 if exist
			if knownVulnSample.SHA1 != "" && knownVulnSample.SHA1 != "-" {
				LoadedHashes.SHA1Sums = append(LoadedHashes.SHA1Sums, knownVulnSample.SHA1)
			}
			// Append SHA256 if exist
			if knownVulnSample.SHA256 != "" && knownVulnSample.SHA256 != "-" {
				LoadedHashes.SHA256Sums = append(LoadedHashes.SHA256Sums, knownVulnSample.SHA256)
			}
		}
	}

	logger.Logf("    |--> %d MD5 hashes", len(LoadedHashes.MD5Sums))
	logger.Logf("    |--> %d SHA1 hashes", len(LoadedHashes.SHA1Sums))
	logger.Logf("    |--> %d SHA256 hashes", len(LoadedHashes.SHA256Sums))

	return nil
}

// MatchHash() will return the matching loldrivers.Driver for a given hash or else will return an error
func MatchHash(hash string) (match *Driver) {
	for _, driver := range LoadedDrivers {
		for _, knownVulnSample := range driver.KnownVulnerableSamples {
			switch len(hash) {
			case 32:
				if knownVulnSample.MD5 == hash {
					return &driver
				}
			case 40:
				if knownVulnSample.SHA1 == hash {
					return &driver
				}
			case 64:
				if knownVulnSample.SHA256 == hash {
					return &driver
				}
			}
		}
	}

	return nil
}
