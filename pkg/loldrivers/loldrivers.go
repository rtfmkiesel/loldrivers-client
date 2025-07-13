package loldrivers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
)

// Will load the drivers based on the selected mode (online, local (requires filePath), internal)
func LoadDrivers(mode string, filePath string) error {
	var (
		jsonData []byte
		err      error
	)

	switch mode {
	case "online":
		jsonData, err = downloadNewestData()
	case "local":
		jsonData, err = os.ReadFile(filePath) // #nosec G304
	case "internal":
		jsonData = internalDrivers
	default:
		return fmt.Errorf("invalid mode: %s", mode)
	}

	if err != nil {
		logger.Error(err)
		logger.Debug("Falling back to internal drivers")
		jsonData = internalDrivers
	}

	if err := loadJsonIntoHashmaps(jsonData); err != nil {
		return err
	}

	return nil
}

// Will return true and a pointer to the matching driver, else return false and nil
func MatchHash(hash string) (bool, *Driver) {
	switch len(hash) {
	case 32:
		if driver := md5Sums[hash]; driver != nil {
			return true, driver
		}
	case 40:
		if driver := sha1Sums[hash]; driver != nil {
			return true, driver
		}
	case 64:
		if driver := sha2Sums[hash]; driver != nil {
			return true, driver
		}
	}

	return false, nil
}

// Downloads the newset driver set from the loldrivers.io API
func downloadNewestData() ([]byte, error) {
	logger.Debug("Downloading the newest data set")

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://www.loldrivers.io/api/drivers.json", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "LOLDrivers-client")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
