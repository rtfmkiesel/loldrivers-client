// Package filesystem handles filesystem operations
package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"loldrivers-client/pkg/logger"
)

var (
	// Default driver paths for Windows 10 and Windows 11
	DriverPaths = []string{"C:\\Windows\\System32\\drivers", "C:\\Windows\\System32\\DriverStore\\FileRepository"}
)

// FileExists() will return true if a given file exists
func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// FileRead() will return the contents of a file as bytes
func FileRead(filepath string) (contentBytes []byte, err error) {
	// Check if file exists
	if !FileExists(filepath) {
		return nil, fmt.Errorf("file '%s' does not exist", filepath)
	}

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s'", filepath)
	}
	defer file.Close()

	// Read file
	contentBytes, err = io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s'", filepath)
	}

	return contentBytes, nil
}

// FilesInFolder() will return a []string of filepaths from a given folder (recursively)
//
// sizeLimit in MB (ex: 5)
func FilesInFolder(path string, sizeLimit int64) (files []string, err error) {
	logger.Verbose(fmt.Sprintf("[*] Searching for files in %s", path))

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

		// Append to slice
		files = append(files, path)
		return nil
	})

	if err != nil {
		return files, err
	}

	logger.Verbose(fmt.Sprintf("[+] Found %d files in %s", len(files), path))
	return files, nil
}

// IsAdmin() will return true is the current binary runs under administrative privileges
func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
