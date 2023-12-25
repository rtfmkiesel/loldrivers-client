package checksums

import (
	"sync"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
	"github.com/rtfmkiesel/loldrivers-client/pkg/loldrivers"
	"github.com/rtfmkiesel/loldrivers-client/pkg/result"
)

// checksums.CalcRunner() is used as a go func for calculating and comparing file checksums
// from a chanJobs. If a calculated checksum matches a loaded checksum a result.Result will be sent to chanResults
func CalcRunner(wg *sync.WaitGroup, chanJobs <-chan string, chanResults chan<- result.Result) {
	defer wg.Done()

	for job := range chanJobs {
		// SHA256
		sha256, err := calcSHA256(job)
		if err != nil {
			logger.Error(err)
		} else if driver := loldrivers.MatchHash(sha256); driver != nil {
			chanResults <- result.Result{
				Filepath: job,
				Checksum: sha256,
				Driver:   *driver,
			}

			continue
		}

		// SHA1
		sha1, err := calcSHA1(job)
		if err != nil {
			logger.Error(err)
		} else if driver := loldrivers.MatchHash(sha1); driver != nil {
			chanResults <- result.Result{
				Filepath: job,
				Checksum: sha1,
				Driver:   *driver,
			}

			continue
		}

		// MD5
		md5, err := calcMD5(job)
		if err != nil {
			logger.Error(err)
		} else if driver := loldrivers.MatchHash(md5); driver != nil {
			chanResults <- result.Result{
				Filepath: job,
				Checksum: md5,
				Driver:   *driver,
			}

			continue
		}
	}
}
