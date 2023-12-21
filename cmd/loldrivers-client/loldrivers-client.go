//go:build windows

package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/rtfmkiesel/loldrivers-client/pkg/checksums"
	"github.com/rtfmkiesel/loldrivers-client/pkg/filesystem"
	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
	"github.com/rtfmkiesel/loldrivers-client/pkg/loldrivers"
	"github.com/rtfmkiesel/loldrivers-client/pkg/output"
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
	var flagWorkers int
	flag.StringVar(&flagMode, "m", "online", "")
	flag.StringVar(&flagMode, "mode", "online", "")
	flag.StringVar(&flagDir, "d", "", "")
	flag.StringVar(&flagDir, "scan-dir", "", "")
	flag.Int64Var(&flagFileLimit, "l", 10, "")
	flag.Int64Var(&flagFileLimit, "scan-limit", 10, "")
	flag.StringVar(&flagLocalFile, "f", "", "")
	flag.StringVar(&flagLocalFile, "driver-file", "", "")
	flag.BoolVar(&flagSilent, "s", false, "")
	flag.BoolVar(&flagSilent, "silent", false, "")
	flag.BoolVar(&flagJSON, "j", false, "")
	flag.BoolVar(&flagJSON, "json", false, "")
	flag.IntVar(&flagWorkers, "w", 20, "")
	flag.IntVar(&flagWorkers, "workers", 20, "")
	flag.Usage = func() {
		logger.Banner()
		fmt.Println(`Usage:
  LOLDrivers-client.exe [OPTIONS]

Options:
  -m, --mode            Operating Mode {online, local, internal}
                            online = Download the newest driver set (default)
                            local = Use a local drivers.json file (requires '-f')
                            internal = Use the built-in driver set (can be outdated)

  -d, --scan-dir        Directory to scan for drivers (default: Windows driver folders)
                        Files which cannot be opened or read will be silently ignored
  -l, --scan-limit      Size limit for files to scan in MB (default: 10)
                        Be aware, higher values greatly increase runtime & CPU usage

  -f, --driver-file     File path to 'drivers.json', when running with '-m local'

  -s, --silent          Will only output found files for easy parsing (default: false)
  -j, --json            Format output as JSON (default: false)

  -w, --workers         Number of "threads" to spawn (default: 20)
  -h, --help            Shows this text
	`)
	}
	flag.Parse()

	// Only one output style
	if flagSilent && flagJSON {
		logger.CatchCrit(fmt.Errorf("only use '-s' or '-j', not both"))
	} else if flagSilent {
		output.Mode = "silent"
		logger.BeSilent = true
	} else if flagJSON {
		output.Mode = "json"
		logger.BeSilent = true
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
	logger.Log(fmt.Sprintf("    |-- Got %d MD5 hashes", len(driverHashes.MD5Sums)))
	logger.Log(fmt.Sprintf("    |-- Got %d SHA1 hashes", len(driverHashes.SHA1Sums)))
	logger.Log(fmt.Sprintf("    |-- Got %d SHA256 hashes", len(driverHashes.SHA256Sums)))

	// Create the channels and waitgroup for the checksum runners
	chanFiles := make(chan string)
	chanResults := make(chan output.Result)
	wgRunner := new(sync.WaitGroup)
	// Spawn the checksum runners
	for i := 0; i <= flagWorkers; i++ {
		go checksums.Runner(wgRunner, chanFiles, chanResults, driverHashes, drivers)
		wgRunner.Add(1)
	}

	// Create the waitgroup for the output runner
	wgOutput := new(sync.WaitGroup)
	// Spawn the output runner
	go output.Runner(wgOutput, chanResults)
	wgOutput.Add(1)

	// Set the folders which are going to be scanned for files
	var paths []string
	if flagDir == "" {
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
