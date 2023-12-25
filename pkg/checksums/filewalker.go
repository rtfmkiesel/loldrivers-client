package checksums

import (
	"os"
	"path/filepath"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
)

// checksums.FileWalker() will recursively send files from path, who are smaller than sizeLimit, to outputChannel
func FileWalker(path string, sizeLimit int64, outputChannel chan<- string) (err error) {
	logger.Logf("[*] Searching for files in %s", path)

	// Walk over every file in a given folder
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Ignore a file if we get "Access is denied" error
			if os.IsPermission(err) {
				return nil
			}

			// Ignore a file if we get "The system cannot find the file specified" error
			if os.IsNotExist(err) {
				return nil
			}

			return err
		}

		// Skip directories and non regular files
		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		// Skip files that can't be read
		if info.Mode().Perm()&0400 == 0 {
			return nil
		}

		// Skip files larger than the specified size limit
		if info.Size() > sizeLimit*1024*1024 {
			return nil
		}

		// Send to the channel
		outputChannel <- path
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
