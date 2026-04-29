package cmd

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rtfmkiesel/loldrivers-client/internal/checksum"
	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
	"github.com/rtfmkiesel/loldrivers-client/internal/loldrivers"
)

type result struct {
	Hash     string             `json:"-"`
	Filepath string             `json:"filepath"`
	Url      string             `json:"-"`
	Driver   *loldrivers.Driver `json:"loldriver"`
}

func Run() error {
	start := time.Now()

	cfg, err := parseConfig()
	if err != nil {
		return err
	}

	if err = loldrivers.LoadDrivers(cfg.mode, cfg.localFile); err != nil {
		return err
	}

	chanFiles := make(chan string, cfg.workerCount)
	chanResults := make(chan *result)
	wgChecksums := new(sync.WaitGroup)

	for i := 1; i <= cfg.workerCount; i++ {
		wgChecksums.Go(func() {
			for file := range chanFiles {
				for _, hashAlg := range checksum.Algorithms {
					hash, err := hashAlg.Calculate(file)
					if err != nil {
						logger.Errorf("hashcalc: type=%s file=%s: %v", hashAlg.Name, file, err)
						continue
					}

					if matched, driver := loldrivers.FindDriverByHash(hash); matched {
						chanResults <- &result{
							Hash:     hash,
							Filepath: file,
							Url:      "https://loldrivers.io/drivers/" + driver.ID,
							Driver:   driver,
						}
						break // Found a match, skip remaining hashes
					}
				}
			}
		})
	}

	wgOutput := new(sync.WaitGroup)
	wgOutput.Go(func() {
		for result := range chanResults {
			switch cfg.outputmode {
			case "standard":
				logger.Warning("MATCHED %s\n\tFILE: %s\n\tURL: %s", result.Hash, result.Filepath, result.Url)
			case "grep":
				logger.Stdout("file=%s;url=%s", result.Filepath, result.Url)
			case "json":
				jsonBytes, err := json.Marshal(result)
				if err != nil {
					logger.Errorf("json marshal: %v", err)
				}
				logger.Stdout("%s", jsonBytes)
			}
		}
	})

	for _, dir := range cfg.targets {
		logger.Debug("Scanning %s", dir)

		if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				// Ignore permission and not-exist errors
				if os.IsPermission(err) || os.IsNotExist(err) {
					return nil
				}

				return err
			}

			// Skip
			if info.IsDir() || // directories
				!info.Mode().IsRegular() || // non-regular files
				info.Mode().Perm()&0400 == 0 || // unreadable files
				info.Size() > cfg.maxSize { // files over the limit
				return nil
			}

			chanFiles <- path
			return nil
		}); err != nil {
			logger.Error(err)
		}
	}

	close(chanFiles)
	wgChecksums.Wait()
	close(chanResults)
	wgOutput.Wait()

	logger.Debug("Finished in %s", time.Since(start))
	return nil
}
