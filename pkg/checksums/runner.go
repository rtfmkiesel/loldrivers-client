package checksums

import (
	"sync"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
	"github.com/rtfmkiesel/loldrivers-client/pkg/loldrivers"
)

type Result struct {
	Filepath string             `json:"Filepath"`
	Checksum string             `json:"Checksum"`
	Driver   *loldrivers.Driver `json:"LOLDrivers-Entry"`
}

// Is used as a go func for calculating and comparing file checksums
func Runner(wg *sync.WaitGroup, filepaths <-chan string, results chan<- *Result, silenceErrors bool) {
	defer wg.Done()

	for filepath := range filepaths {
		// The order of the hash checks was choosen based on the amount of hashes upon writing this bit
		// SHA1 > SHA2 > MD5

		sha1, err := calcSHA1(filepath)
		if err != nil && !silenceErrors {
			logger.Error(err)
			continue
		} else {
			matched, driver := loldrivers.MatchHash(sha1)
			if matched {
				results <- &Result{
					Filepath: filepath,
					Checksum: sha1,
					Driver:   driver,
				}

				continue // No need to check others as there was a match
			}
		}

		sha256, err := calcSHA256(filepath)
		if err != nil && !silenceErrors {
			logger.Error(err)
			continue
		} else {
			matched, driver := loldrivers.MatchHash(sha256)
			if matched {
				results <- &Result{
					Filepath: filepath,
					Checksum: sha256,
					Driver:   driver,
				}

				continue // No need to check others as there was a match
			}
		}

		md5, err := calcMD5(filepath)
		if err != nil && !silenceErrors {
			logger.Error(err)
		} else {
			matched, driver := loldrivers.MatchHash(md5)
			if matched {
				results <- &Result{
					Filepath: filepath,
					Checksum: md5,
					Driver:   driver,
				}
			}
		}
	}
}
