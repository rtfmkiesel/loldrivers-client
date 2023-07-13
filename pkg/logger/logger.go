// Package logger handles errors as well as the output handling
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	// To control verbose output
	BeVerbose      bool = false
	OutputJSON     bool = false
	OutputParsable bool = false
)

const (
	// Time format for the terminal output
	timeFormat = "2006-01-02T15:04:05"
)

// Struct for the results of the compare runners
type Result struct {
	Filename string
	Checksum string
}

// Log() will print to the terminal
func Log(message string) {
	if !OutputParsable && !OutputJSON {
		fmt.Printf("%s %s\n", time.Now().Format(timeFormat), message)
	}
}

// Verbose() will print verbose messages to the terminal if verbose mode is selected
func Verbose(message string) {
	if BeVerbose {
		fmt.Printf("%s %s\n", time.Now().Format(timeFormat), message)
	}
}

// Catch() will handle errors
func Catch(err error) {
	fmt.Printf("%s ERROR: %s\n", time.Now().Format(timeFormat), err)

}

// CatchCrit() will handle critical errors
func CatchCrit(err error) {
	fmt.Printf("%s CRITICAL: %s\n", time.Now().Format(timeFormat), err)
	os.Exit(1)
}

// Banner() will print the banner
func Banner() {
	if !OutputParsable && !OutputJSON {
		fmt.Printf(`   __    ___  __     _      _                    
  / /   /___\/ /  __| |_ __(_)_   _____ _ __ ___ 
 / /   //  // /  / _  | '__| \ \ / / _ \ '__/ __|
/ /___/ \_// /__| (_| | |  | |\ V /  __/ |  \__ \
\____/\___/\____/\__,_|_|  |_| \_/ \___|_|  |___/
https://www.loldrivers.io | Client by @rtfmkiesel

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
		// Print result
		if OutputParsable {
			fmt.Printf("%s;%s\n", job.Filename, job.Checksum)
		} else if OutputJSON {
			jsonOutput, err := json.Marshal(job)
			if err != nil {
				Catch(err)
				continue
			}
			fmt.Printf("%s\n", string(jsonOutput))
		} else {
			Log(fmt.Sprintf("[+] %s:%s", job.Filename, job.Checksum))
		}
		// Increment counter
		counter++
	}

	if counter == 0 {
		Log("[-] No vulnerable/malicious driver(s) found!")
	} else {
		Log(fmt.Sprintf("[*] Found %d vulnerable/malicious driver(s)!", counter))
	}
}
