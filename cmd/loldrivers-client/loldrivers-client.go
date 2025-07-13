//go:build windows

package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/fatih/color"

	"github.com/rtfmkiesel/loldrivers-client/pkg/checksums"
	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
	"github.com/rtfmkiesel/loldrivers-client/pkg/loldrivers"
	"github.com/rtfmkiesel/loldrivers-client/pkg/options"
)

func main() {
	opt, err := options.Parse()
	if err != nil {
		logger.Fatal(err)
	}

	if err = loldrivers.LoadDrivers(opt.Mode, opt.ModeLocalFilePath); err != nil {
		logger.Fatal(err)
	}

	chanFilepaths := make(chan string)
	chanResults := make(chan *checksums.Result)
	wgChecksums := new(sync.WaitGroup)
	for i := 0; i <= opt.ScanWorkers; i++ {
		go checksums.Runner(wgChecksums, chanFilepaths, chanResults, opt.ScanShowErrors)
		wgChecksums.Add(1)
	}

	wgOutput := new(sync.WaitGroup)
	go func() {
		defer wgOutput.Done()
		for result := range chanResults {
			switch opt.OutputMode {
			case "grep":
				logger.Stdout("%s", result.Filepath)
			case "json":
				jsonOutput, err := json.Marshal(result)
				if err != nil {
					logger.Error(err)
					continue
				}
				logger.Stdout("%s", string(jsonOutput))
			default:
				logger.Custom("VUL", color.FgRed, "%s (https://loldrivers.io/drivers/%s)", result.Filepath, result.Driver.ID)
			}
		}
	}()
	wgOutput.Add(1)

	for _, path := range opt.ScanDirectories {
		if err := checksums.DirectoryWalker(path, opt.ScanSizeLimit, chanFilepaths); err != nil {
			logger.Error(err)
			continue
		}
	}

	close(chanFilepaths)
	wgChecksums.Wait()
	close(chanResults)
	wgOutput.Wait()

	logger.Debug("Finished in %s", time.Since(opt.StartTime))
}
