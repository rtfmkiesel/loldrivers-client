//go:build windows

package main

import (
	"sync"
	"time"

	"github.com/rtfmkiesel/loldrivers-client/pkg/checksums"
	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
	"github.com/rtfmkiesel/loldrivers-client/pkg/loldrivers"
	"github.com/rtfmkiesel/loldrivers-client/pkg/options"
	"github.com/rtfmkiesel/loldrivers-client/pkg/result"
)

func main() {
	// Parse the command line options
	opt, err := options.Parse()
	if err != nil {
		logger.Fatal(err)
	}

	// Load the drivers and their hashes
	if err = loldrivers.LoadDrivers(opt.Mode, opt.LocalDriversPath); err != nil {
		logger.Fatal(err)
	}

	// Set up the checksum runners
	chanFiles := make(chan string)
	chanResults := make(chan result.Result)
	wgRunner := new(sync.WaitGroup)
	for i := 0; i <= opt.Workers; i++ {
		go checksums.CalcRunner(wgRunner, chanFiles, chanResults)
		wgRunner.Add(1)
	}

	// Set up the one output runner
	wgResults := new(sync.WaitGroup)
	go result.OutputRunner(wgResults, chanResults, opt.OutputMode)
	wgResults.Add(1)

	// Get all files from subfolders and send them to the checksum runners via a channel
	for _, path := range opt.ScanDirectories {
		if err := checksums.FileWalker(path, opt.ScanSizeLimit, chanFiles); err != nil {
			logger.Fatal(err)
		}
	}

	close(chanFiles)
	wgRunner.Wait()

	close(chanResults)
	wgResults.Wait()

	logger.Logf("[+] Done, took %s", time.Since(opt.StartTime))
}
