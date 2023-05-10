//go:build windows

package main

import (
	"flag"
	"fmt"
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
	var flagThreads int
	var flagVerbose bool
	flag.StringVar(&flagMode, "m", "online", "")
	flag.StringVar(&flagDir, "d", "", "")
	flag.Int64Var(&flagFileLimit, "l", 10, "")
	flag.StringVar(&flagLocalFile, "f", "", "")
	flag.IntVar(&flagThreads, "t", 20, "")
	flag.BoolVar(&flagVerbose, "v", false, "")
	flag.Usage = func() {
		fmt.Println(`
LOLDrivers-client.exe -m [MODE] [OPTIONS]

Modes:
  online    Download the newest driver set (default)
  local     Use a local drivers.json file (requires '-f')
  internal  Use the built-in driver set (can be outdated, fallback)

Options:
  -d        Directory to scan for drivers (default: Windows Default Driver Folders)
            Files which cannot be opened or read will be silently ignored
  -l        Size limit for files to scan in MB (default: 10)
            Be aware, higher values greatly increase runtime & CPU usage
  -f        File path to 'drivers.json'
            Only needed with '-m local'
  -t        Number of threads to spawn (default: 20)
  -v        Print verbose messages (default: false)
  -h        Shows this text
	`)
	}
	flag.Parse()

	if flagVerbose {
		// Set the value in the logger module
		logger.BeVerbose = true
	}

	// ASCII L0VE
	logger.Banner()
	logger.Log("[*] Started")

	// To store the drivers in
	var drivers []loldrivers.Driver

	// Load the drivers based on the selected mode
	logger.Verbose(fmt.Sprintf("[*] Mode '%s'", flagMode))

	switch flagMode {
	case "online":
		// Default, download from the web
		// Download drivers
		var err error
		drivers, err = loldrivers.Get()
		if err != nil {
			// There was a parsing error
			logger.Catch(err)
			logger.Log("[!] Got an error while parsing online data. Falling back to internal data set")
			drivers, err = loldrivers.Parse(loldrivers.InternalDrivers)
			if err != nil {
				logger.CatchCrit(err)
			}
		}

	case "local":
		// User wants to use a local file
		if flagLocalFile == "" {
			logger.CatchCrit(fmt.Errorf("mode 'local' requires '-f'"))
		}

		// Read file
		jsonBytes, err := filesystem.FileRead(flagLocalFile)
		if err != nil {
			logger.CatchCrit(err)
		}

		// Parse file
		drivers, err = loldrivers.Parse(jsonBytes)
		if err != nil {
			// There was a parsing error
			logger.Catch(err)
			logger.Log("[!] Got an error while parsing local file. Falling back to internal data set")
			drivers, err = loldrivers.Parse(loldrivers.InternalDrivers)
			if err != nil {
				logger.CatchCrit(err)
			}
		}

	case "internal":
		// User wants to use internal data set
		// Parse bytes
		var err error
		drivers, err = loldrivers.Parse(loldrivers.InternalDrivers)
		if err != nil {
			logger.CatchCrit(err)
		}

	default:
		logger.CatchCrit(fmt.Errorf("invalid mode '%s'", flagMode))
	}

	logger.Verbose(fmt.Sprintf("[*] Loaded %d drivers", len(drivers)))

	// To store all driver checksums
	driverHashes := loldrivers.DriverHashes{}
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

	logger.Verbose(fmt.Sprintf("[*] Got %d MD5 hashes", len(driverHashes.MD5Sums)))
	logger.Verbose(fmt.Sprintf("[*] Got %d SHA1 hashes", len(driverHashes.SHA1Sums)))
	logger.Verbose(fmt.Sprintf("[*] Got %d SHA256 hashes", len(driverHashes.SHA256Sums)))

	// Create the channels and waitgroup for the checksum runners
	chanFiles := make(chan string)
	chanChecksums := make(chan checksums.Sums)
	wgChecksums := new(sync.WaitGroup)
	// Spawn the checksum runners
	for i := 0; i <= flagThreads; i++ {
		go checksums.CalcRunner(wgChecksums, chanFiles, chanChecksums)
		wgChecksums.Add(1)
	}

	// Create the channels and waitgroup for the compare runners
	chanResults := make(chan logger.Result)
	wgCompare := new(sync.WaitGroup)
	// Spawn the compare runners
	for i := 0; i <= flagThreads; i++ {
		go checksums.CompareRunner(wgCompare, chanChecksums, chanResults, driverHashes)
		wgCompare.Add(1)
	}

	// Create the waitgroup for the output runner
	wgOutput := new(sync.WaitGroup)
	// Spawn the output runner
	go logger.OutputRunner(wgOutput, chanResults)
	wgOutput.Add(1)

	// Scan for drivers
	//
	// A tiny amount of time could be saved if only ".sys" files were checked.
	// I tried and the time saved is not worth it since malicious drivers could have other file extensions.
	if flagDir == "" {
		// User did not specify a path, use the default Windows paths

		// Since scanning the default folders requires admin privileges, check here
		if !filesystem.IsAdmin() {
			logger.CatchCrit(fmt.Errorf("not running with administrative privileges"))
		}

		for _, path := range filesystem.DriverPaths {
			// Get all files
			files, err := filesystem.FilesInFolder(path, flagFileLimit)
			if err != nil {
				logger.CatchCrit(err)
			}

			// Add jobs to the channel
			for _, file := range files {
				chanFiles <- file
			}
		}
	} else {
		// Scan the user specified folder for drivers
		// Get all files
		files, err := filesystem.FilesInFolder(flagDir, flagFileLimit)
		if err != nil {
			logger.CatchCrit(err)
		}

		// Add jobs to the channel
		for _, file := range files {
			chanFiles <- file
		}
	}

	// Close the channel to start the checksum runners
	close(chanFiles)
	// Wait here until all checksums are calculated
	logger.Verbose("[*] Waiting for the checksum runners to complete")
	wgChecksums.Wait()

	// Close the channel to start the compare runners
	close(chanChecksums)
	// Wait here until all compares have been done
	logger.Verbose("[*] Waiting for the compare runners to complete")
	wgCompare.Wait()

	// Close the results channel to process the results
	close(chanResults)
	// Wait until all results have been processed
	wgOutput.Wait()

	logger.Verbose(fmt.Sprintf("[*] Done, took %s\n", time.Since(startTime)))
}
