package result

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
	"github.com/rtfmkiesel/loldrivers-client/pkg/loldrivers"
)

// Struct for the results of the compare runners
type Result struct {
	Filepath string
	Checksum string
	Driver   loldrivers.Driver
}

// OutputRunner() is used as a go func to display the results
func OutputRunner(wg *sync.WaitGroup, chanResults <-chan Result, mode string) {
	defer wg.Done()

	// To count how many results we got
	counter := 0

	for result := range chanResults {
		switch mode {
		case "silent":
			fmt.Fprintf(os.Stdout, "%s\n", result.Filepath)
		case "json":
			jsonOutput, err := json.Marshal(result)
			if err != nil {
				logger.Fatal(err)
			}
			fmt.Fprintf(os.Stdout, "%s\n", string(jsonOutput))
		default:
			fmt.Fprintf(os.Stdout, "[!] Found %s\n", result.Driver.Category)
			fmt.Fprintf(os.Stdout, "    |--> Path: %s\n", result.Filepath)
			//fmt.Fprintf(os.Stdout, "    |--> Checksum: %s\n", result.Checksum)
			fmt.Fprintf(os.Stdout, "    |--> Link: https://loldrivers.io/drivers/%s\n", result.Driver.ID)
		}

		counter++
	}

	if counter == 0 {
		logger.Log("[-] No vulnerable or malicious driver(s) found!")
	} else {
		logger.Logf("[+] Found a total of %d vulnerable or malicious driver(s)!", counter)
	}
}
