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

	chanFilepaths := make(chan string, opt.ScanWorkers*2)
	chanResults := make(chan *checksums.Result)
	wgChecksums := new(sync.WaitGroup)

	for i := 1; i <= opt.ScanWorkers; i++ {
		wgChecksums.Go(func() {
			checksums.Runner(chanFilepaths, chanResults, opt.ScanShowErrors)
		})
	}

	wgOutput := new(sync.WaitGroup)
	wgOutput.Go(func() {
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
	})

	for _, path := range opt.ScanDirectories {
		if err := checksums.DirectoryWalker(path, opt.ScanSizeLimit, chanFilepaths); err != nil {
			logger.Error(err)
		}
	}

	close(chanFilepaths)
	wgChecksums.Wait()
	close(chanResults)
	wgOutput.Wait()

	logger.Debug("Finished in %s", time.Since(opt.StartTime))
}
