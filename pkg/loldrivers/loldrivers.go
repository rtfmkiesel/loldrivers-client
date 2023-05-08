// Package loldrivers handles the JSON data from loldrivers.io
package loldrivers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"loldrivers-client/pkg/logger"
	"net/http"
)

var (
	// Embed a driver.json during build for use with -m 'internal'
	//go:embed drivers.json
	InternalDrivers []byte
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
	Commands        UnmarshalCommands `json:"Commands,omitempty"`
	Resources       []string          `json:"Resources,omitempty"`
	Acknowledgement struct {
		Person StringOrStringArray `json:"Person"`
		Handle string              `json:"Handle"`
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

// Struct that is used during unmarshalling of the JSON data
// since sometimes "Commands" will be either a single string or a "Command" struct
type UnmarshalCommands struct {
	Value []Command
	Set   bool
}

// Struct that is used during unmarshalling of the JSON data
// since sometimes a key can be either a single string or an array of strings
type StringOrStringArray struct {
	Value []string
	Set   bool
}

// Struct to store the driver hashes from loldrivers.io
type DriverHashes struct {
	MD5Sums    []string
	SHA1Sums   []string
	SHA256Sums []string
}

// The UnmarshalJSON method on UnmarshalCommands will parse the JSON
// as eiter a "Command" struct or a single string (into a "Command" struct)
func (s *UnmarshalCommands) UnmarshalJSON(b []byte) error {
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
func (s *StringOrStringArray) UnmarshalJSON(b []byte) error {
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

// Get() will download the current loldrivers.io driver data set
func Get() (drivers []Driver, err error) {
	logger.Verbose(fmt.Sprintf("[*] Downloading the newest drivers from %s", apiURL))

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
	logger.Verbose("[+] Download successful")

	// Read the bode into []byte
	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse the data
	drivers, err = Parse(jsonBytes)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}

// Parse() will return a slice of drivers from JSON input bytes
func Parse(jsonBytes []byte) (drivers []Driver, err error) {
	logger.Verbose("[*] Parsing JSON data into struct")

	// Unmarshal JSON data
	if err := json.Unmarshal(jsonBytes, &drivers); err != nil {
		return nil, err
	}

	logger.Verbose("[+] Parsing successful")
	return drivers, nil
}
