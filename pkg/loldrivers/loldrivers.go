// Package loldrivers handles the JSON data from loldrivers.io
package loldrivers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"loldrivers-client/pkg/filesystem"
	"loldrivers-client/pkg/logger"
	"net/http"
)

var (
	// Embed a driver.json during build for use with -m 'internal'
	//go:embed drivers.json
	internalDrivers []byte
)

const (
	// Download link to the 'drivers.json' file
	apiURL = "https://www.loldrivers.io/api/drivers.json"
)

// Struct for a single driver from loldrivers.io
//
// Based on the the JSON spec from
// https://github.com/magicsword-io/LOLDrivers/blob/validate/bin/spec/drivers.spec.json
type Driver struct {
	ID              string            `json:"Id"`
	Author          string            `json:"Author"`
	Created         string            `json:"Created"`
	MitreID         string            `json:"MitreID"`
	Category        string            `json:"Category"`
	Verified        string            `json:"Verified"`
	Commands        unmarshalCommands `json:"Commands,omitempty"`
	Resources       []string          `json:"Resources,omitempty"`
	Acknowledgement struct {
		Person unmarshalStringOrStringArray `json:"Person"`
		Handle string                       `json:"Handle"`
	} `json:"Acknowledgement,omitempty"`
	Detection []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"Detection,omitempty"`
	KnownVulnerableSamples []struct {
		Filename string `json:"Filename"`
		MD5      string `json:"MD5,omitempty"`
		SHA1     string `json:"SHA1,omitempty"`
		SHA256   string `json:"SHA256,omitempty"`
	} `json:"KnownVulnerableSamples,omitempty"`
	Tags []string `json:"Tags"`
}

// 'Command' struct for a driver from loldrivers.io
//
// Based on the the JSON spec from
// https://github.com/magicsword-io/LOLDrivers/blob/validate/bin/spec/drivers.spec.json
type Command struct {
	Command         string `json:"Command"`
	Description     string `json:"Description"`
	Usecase         string `json:"Usecase"`
	Privileges      string `json:"Privileges"`
	OperatingSystem string `json:"OperatingSystem"`
}

// Struct to store the driver hashes from loldrivers.io
type DriverHashes struct {
	MD5Sums    []string
	SHA1Sums   []string
	SHA256Sums []string
}

// Struct that is used during unmarshalling of the "Commands" JSON data
// since sometimes it'll be either a single string or a "Command" struct
type unmarshalCommands struct {
	Value []Command
	Set   bool
}

// Struct that is used during unmarshalling of various the JSON data
// since sometimes a key can be either a single string or an array of strings
type unmarshalStringOrStringArray struct {
	Value []string
	Set   bool
}

// The UnmarshalJSON method on UnmarshalCommands will parse the JSON
// as eiter a "Command" struct or a single string (into a "Command" struct)
func (s *unmarshalCommands) UnmarshalJSON(b []byte) error {
	var strVal string
	var cmdVal Command
	// Try to unmarshal into a string first
	err := json.Unmarshal(b, &strVal)
	if err == nil {
		// No error, set string value, leave rest empty
		cmdVal = Command{
			Command: strVal,
		}
	} else {
		// Try to unmarshall into a "Command" struct
		err = json.Unmarshal(b, &cmdVal)
		if err != nil {
			// Both unmarshall were unsuccessful
			return err
		}
	}
	// Set the value of s to the unmarshalled value
	s.Value = append(s.Value, cmdVal)
	s.Set = true
	return nil
}

// The UnmarshalJSON method will parse the JSON as either a single string
// or an array of strings into a slice of strings
func (s *unmarshalStringOrStringArray) UnmarshalJSON(b []byte) error {
	var strVal string
	var arrVal []string
	// Try to unmarshal into a single string first
	err := json.Unmarshal(b, &strVal)
	if err == nil {
		// No error, create a array with a single value
		arrVal = []string{strVal}
	} else {
		// Try to unmarshall into a slice of strings
		err = json.Unmarshal(b, &arrVal)
		if err != nil {
			// Both unmarshall were unsuccessful
			return err
		}
	}
	// Set the value of s to the unmarshalled value which will always be a slice
	s.Value = arrVal
	s.Set = true
	return nil
}

// LoadDrivers() will load the drivers based on the selected mode
// and return a slice of loldrivers.Driver
//
// mode = online, local, internal
//
// filepath = path to JSON file if mode == local, else ""
func LoadDrivers(mode string, filepath string) (drivers []Driver, err error) {
	// Load the drivers based on the selected mode
	logger.Log(fmt.Sprintf("[*] Loading drivers with mode '%s'", mode))

	switch mode {
	case "online":
		// Default, download from the web
		// Download drivers
		drivers, err = download()
		if err != nil {
			// There was a parsing error
			logger.Catch(err)
			logger.Log("[!] Got an error while parsing online data. Falling back to internal data set")
			drivers, err = parse(internalDrivers)
			if err != nil {
				return drivers, err
			}
		}

	case "local":
		// User wants to use a local file
		if filepath == "" {
			logger.CatchCrit(fmt.Errorf("mode 'local' requires '-f'"))
		}

		// Read file
		jsonBytes, err := filesystem.FileRead(filepath)
		if err != nil {
			return drivers, err
		}

		// Parse file
		drivers, err = parse(jsonBytes)
		if err != nil {
			// There was a parsing error
			logger.Catch(err)
			logger.Log("[!] Got an error while parsing local file. Falling back to internal data set")
			drivers, err = parse(internalDrivers)
			if err != nil {
				return drivers, err
			}
		}

	case "internal":
		// User wants to use internal data set
		// Parse bytes
		drivers, err = parse(internalDrivers)
		if err != nil {
			return drivers, err
		}

	default:
		logger.CatchCrit(fmt.Errorf("invalid mode '%s'", mode))
	}

	return drivers, nil
}

// GetHashes() will return loldrivers.DriverHashes containing all
// MD5, SHA1 and SHA256 from a slice of loldrivers.Driver. Empty values or '-' will be ignored
func GetHashes(drivers []Driver) (driverHashes DriverHashes) {
	// Get all checksums from the loaded drivers
	for _, driver := range drivers {
		for _, knownVulnSample := range driver.KnownVulnerableSamples {
			// Append MD5 if exist
			if knownVulnSample.MD5 != "" && knownVulnSample.MD5 != "-" {
				driverHashes.MD5Sums = append(driverHashes.MD5Sums, knownVulnSample.MD5)
			}
			// Append SHA1 if exist
			if knownVulnSample.SHA1 != "" && knownVulnSample.SHA1 != "-" {
				driverHashes.SHA1Sums = append(driverHashes.SHA1Sums, knownVulnSample.SHA1)
			}
			// Append SHA256 if exist
			if knownVulnSample.SHA256 != "" && knownVulnSample.SHA256 != "-" {
				driverHashes.SHA256Sums = append(driverHashes.SHA256Sums, knownVulnSample.SHA256)
			}
		}
	}

	return driverHashes
}

// MatchHash() will return the matching loldrivers.Driver for a given hash or else will return an error
func MatchHash(hash string, drivers []Driver) (match Driver, err error) {
	// Get all checksums from the loaded drivers
	for _, driver := range drivers {
		for _, knownVulnSample := range driver.KnownVulnerableSamples {
			if knownVulnSample.MD5 == hash {
				return driver, nil
			}
			if knownVulnSample.SHA1 == hash {
				return driver, nil
			}
			if knownVulnSample.SHA256 == hash {
				return driver, nil
			}
		}
	}

	return match, fmt.Errorf("no match found")
}

// download() will download the current loldrivers.io data set
func download() (drivers []Driver, err error) {
	logger.Log("[*] Downloading the newest drivers")

	// Setup HTTP client
	client := &http.Client{}

	// Build request
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Make the request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	logger.Log("[+] Download successful")

	// Read the bode into []byte
	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse the data
	drivers, err = parse(jsonBytes)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}

// parse() will return a slice of loldrivers.Drivers from JSON input bytes
func parse(jsonBytes []byte) (drivers []Driver, err error) {
	// Unmarshal JSON data
	if err := json.Unmarshal(jsonBytes, &drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}
