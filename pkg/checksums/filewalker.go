package checksums

import (
	"os"
	"path/filepath"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
)

// Will recursively walk and send filepaths from root, who are smaller than sizeLimit to filepaths
func DirectoryWalker(root string, sizeLimit int, filepaths chan<- string) (err error) {
	logger.Info("Checking %s", root)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
		if info.Size() > int64(sizeLimit)*1024*1024 {
			return nil
		}

		filepaths <- path
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
