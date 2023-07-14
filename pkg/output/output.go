// package output handles the printing of the results to the terminal
package output

import (
	"encoding/json"
	"fmt"
	"loldrivers-client/pkg/logger"
	"loldrivers-client/pkg/loldrivers"
	"sync"
)

var (
	Mode = "default"
)

// Struct for the results of the compare runners
type Result struct {
	Filepath string
	Checksum string
	Driver   loldrivers.Driver
}

// Runner() is used as a go func to display the results
func Runner(wg *sync.WaitGroup, chanResults <-chan Result) {
	defer wg.Done()

	// To count how many results we got
	counter := 0

	// For each result
	for result := range chanResults {
		// Print result based on output
		switch Mode {
		case "silent":
			fmt.Printf("%s\n", result.Filepath)
		case "json":
			jsonOutput, err := json.Marshal(result)
			if err != nil {
				logger.CatchCrit(err)
			}
			fmt.Printf("%s\n", string(jsonOutput))
		default:
			fmt.Printf("[!] MATCH: %s\n", result.Filepath)
			fmt.Printf("    |-- Category: %s\n", result.Driver.Category)
			fmt.Printf("    |-- Checksum: %s\n", result.Checksum)
			fmt.Printf("    |-- Link: https://loldrivers.io/drivers/%s\n", result.Driver.ID)
		}

		counter++
	}

	if counter == 0 {
		logger.Log("[-] No vulnerable or malicious driver(s) found!")
	} else {
		logger.Log(fmt.Sprintf("[+] Found a total of %d vulnerable or malicious driver(s)!", counter))
	}
}
