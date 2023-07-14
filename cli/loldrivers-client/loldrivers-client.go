//go:build windows

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"loldrivers-client/pkg/checksums"
	"loldrivers-client/pkg/filesystem"
	"loldrivers-client/pkg/logger"
	"loldrivers-client/pkg/loldrivers"
)

func main() {
	// To track execution time
	startTime := time.Now()

	// Setup & parse command line arguments
	var flagMode string
	var flagDir string
	var flagFileLimit int64
	var flagLocalFile string
	var flagSilent bool
	var flagJSON bool
	var flagThreads int
	flag.StringVar(&flagMode, "m", "online", "")
	flag.StringVar(&flagDir, "d", "", "")
	flag.Int64Var(&flagFileLimit, "l", 10, "")
	flag.StringVar(&flagLocalFile, "f", "", "")
	flag.BoolVar(&flagSilent, "s", false, "")
	flag.BoolVar(&flagJSON, "j", false, "")
	flag.IntVar(&flagThreads, "t", 20, "")
	flag.Usage = func() {
		logger.Banner()
		fmt.Println(`LOLDrivers-client.exe -m [MODE] [OPTIONS]

Modes:
  online    Download the newest driver set (default)
  local     Use a local drivers.json file (requires '-f')
  internal  Use the built-in driver set (can be outdated, fallback)

Options:
  -d        Directory to scan for drivers (default: Windows driver folders)
            Files which cannot be opened or read will be silently ignored
  -l        Size limit for files to scan in MB (default: 10)
            Be aware, higher values greatly increase runtime & CPU usage

  -f        File path to 'drivers.json'
            Only needed with '-m local'

  -s        Silent (parsable) output (default: false)
  -j        JSON output (default: false)

  -t        Number of threads to spawn (default: 20)
  -h        Shows this text
	`)
	}
	flag.Parse()

	// Only one output style
	if flagSilent && flagJSON {
		logger.CatchCrit(fmt.Errorf("only use '-s' or '-j', not both"))
	} else if flagSilent {
		logger.OutputMode = "silent"
	} else if flagJSON {
		logger.OutputMode = "json"
	}

	// ASCII L0VE
	logger.Banner()
	// Only run on Windows
	if runtime.GOOS != "windows" {
		logger.CatchCrit(fmt.Errorf("this client was made for Windows only"))
	}

	// Load the drivers
	drivers, err := loldrivers.LoadDrivers(flagMode, flagLocalFile)
	if err != nil {
		logger.CatchCrit(err)
	}
	logger.Log(fmt.Sprintf("[+] Loaded %d drivers", len(drivers)))

	// Get all hashes from the loaded drivers
	driverHashes := loldrivers.GetHashes(drivers)
	logger.Log(fmt.Sprintf("[+] Got %d MD5 hashes", len(driverHashes.MD5Sums)))
	logger.Log(fmt.Sprintf("[+] Got %d SHA1 hashes", len(driverHashes.SHA1Sums)))
	logger.Log(fmt.Sprintf("[+] Got %d SHA256 hashes", len(driverHashes.SHA256Sums)))

	// Create the channels and waitgroup for the checksum runners
	chanFiles := make(chan string)
	chanResults := make(chan logger.Result)
	wgRunner := new(sync.WaitGroup)
	// Spawn the checksum runners
	for i := 0; i <= flagThreads; i++ {
		go checksums.Runner(wgRunner, chanFiles, chanResults, driverHashes)
		wgRunner.Add(1)
	}

	// Create the waitgroup for the output runner
	wgOutput := new(sync.WaitGroup)
	// Spawn the output runner
	go logger.OutputRunner(wgOutput, chanResults)
	wgOutput.Add(1)

	// Set the folders which are going to be scanned for files
	var paths []string
	if flagDir == "" {
		// Since scanning the default folders requires admin privileges, check here
		if _, err := os.Open("\\\\.\\PHYSICALDRIVE0"); err != nil {
			logger.CatchCrit(fmt.Errorf("not running with administrative privileges"))
		}

		// User did not specify a path with '-d', use the default Windows paths
		paths = append(paths, "C:\\Windows\\System32\\drivers")
		paths = append(paths, "C:\\Windows\\System32\\DriverStore\\FileRepository")
		paths = append(paths, "C:\\WINDOWS\\inf")
	} else {
		// User specified a custom folder to scan
		paths = append(paths, flagDir)
	}

	// Get all files from subfolders and send them to the checksum runners via a channel
	for _, path := range paths {
		if err := filesystem.FileWalker(path, flagFileLimit, chanFiles); err != nil {
			logger.CatchCrit(err)
		}
	}

	// Close the channel to start the checksum runners
	close(chanFiles)
	// Wait here until all checksums are calculated and compared
	wgRunner.Wait()

	// Close the results channel to process the results
	close(chanResults)
	// Wait until all results have been processed
	wgOutput.Wait()

	logger.Log(fmt.Sprintf("[+] Done, took %s\n", time.Since(startTime)))
}
