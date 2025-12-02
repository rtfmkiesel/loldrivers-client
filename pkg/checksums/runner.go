package checksums

import (
	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
	"github.com/rtfmkiesel/loldrivers-client/pkg/loldrivers"
)

type Result struct {
	Filepath string             `json:"Filepath"`
	Checksum string             `json:"Checksum"`
	Driver   *loldrivers.Driver `json:"LOLDrivers-Entry"`
}

// Is used as a go func for calculating and comparing file checksums
func Runner(filepaths <-chan string, results chan<- *Result, silenceErrors bool) {
	hashFuncs := []struct {
		name string
		fn   func(string) (string, error)
	}{
		{"SHA1", calcSHA1},
		{"SHA256", calcSHA256},
		{"MD5", calcMD5},
	}

	for filepath := range filepaths {
		for _, hash := range hashFuncs {
			checksum, err := hash.fn(filepath)
			if err != nil {
				if !silenceErrors {
					logger.Errorf("Error calculating %s for %s: %v", hash.name, filepath, err)
				}
				continue
			}

			if matched, driver := loldrivers.MatchHash(checksum); matched {
				results <- &Result{
					Filepath: filepath,
					Checksum: checksum,
					Driver:   driver,
				}
				break // Found a match, skip remaining hashes
			}
		}
	}
}
