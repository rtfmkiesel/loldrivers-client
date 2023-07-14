// Package logger handles errors as well as the output handling
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	OutputMode = "default"
)

// Struct for the results of the compare runners
type Result struct {
	Filepath string
	Checksum string
}

// Log() will print to the terminal
func Log(message string) {
	// Only log if mode is default since in
	// JSON or silent mode we do not want this output
	if OutputMode == "default" {
		fmt.Printf("%s\n", message)
	}
}

// Catch() will handle errors
func Catch(err error) {
	fmt.Printf("ERROR: %s\n", err)
}

// CatchCrit() will handle critical errors
func CatchCrit(err error) {
	fmt.Printf("CRITICAL ERROR: %s\n", err)
	os.Exit(1)
}

// Banner() will print the banner
func Banner() {
	if OutputMode == "default" || OutputMode == "verbose" {
		fmt.Printf(`
    ╒═══════════════════════════╕
    |     LOLDrivers-client     |
    |   https://loldrivers.io   |
    |                           |
    | by @rtfmkiesel/mkiesel.ch |
    ╘═══════════════════════════╛

`)
	}
}

// OutputRunner() is used as a go func to display the results
func OutputRunner(wg *sync.WaitGroup, chanJobs <-chan Result) {
	defer wg.Done()

	// To count how many results we got
	counter := 0

	// For each job
	for job := range chanJobs {
		// Print result based on output
		switch OutputMode {
		case "silent":
			fmt.Printf("%s\n", job.Filepath)
		case "json":
			jsonOutput, err := json.Marshal(job)
			if err != nil {
				CatchCrit(err)
			}
			fmt.Printf("%s\n", string(jsonOutput))
		default:
			fmt.Printf("[!] FOUND %s\n", job.Filepath)
		}

		counter++
	}

	if counter == 0 {
		Log("[-] No vulnerable or malicious driver(s) found!")
	} else {
		Log(fmt.Sprintf("[+] Found a total of %d vulnerable or malicious driver(s)!", counter))
	}
}
