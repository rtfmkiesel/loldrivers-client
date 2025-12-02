package checksums

import (
	"os"
	"path/filepath"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
)

// Will recursively walk and send filepaths from root, who are smaller than sizeLimit to filepaths
func DirectoryWalker(root string, sizeLimit int, filepaths chan<- string) error {
	logger.Debug("Checking %s", root)

	sizeLimitBytes := int64(sizeLimit) * 1024 * 1024

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Ignore permission and not-exist errors
			if os.IsPermission(err) || os.IsNotExist(err) {
				return nil
			}

			return err
		}

		// Skip directories, non-regular files, unreadable files, and oversized files
		if info.IsDir() ||
			!info.Mode().IsRegular() ||
			info.Mode().Perm()&0400 == 0 ||
			info.Size() > sizeLimitBytes {
			return nil
		}

		filepaths <- path
		return nil
	})
}
