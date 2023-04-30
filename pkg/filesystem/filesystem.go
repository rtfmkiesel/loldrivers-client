// Package filesystem handles filesystem operations
package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

// FilesInFolder() will return a []string of files in a folder (recursively)
func FilesInFolder(path string) (files []string, err error) {
	logger.Verbose(fmt.Sprintf("[*] Searching for files in %s", path))

	// Walk over every file in a given folder
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only append if it's not a directory
		if !info.IsDir() {
			// Append to slice
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return files, err
	}

	logger.Verbose(fmt.Sprintf("[+] Found %d files in %s", len(files), path))
	return files, nil
}

// FilesInFolderExt() will return a []string of files in a folder (recursively)
//
// A file extension may be specified for filtering, ex: '.txt'.
// If not needed, set 'ext' to an empty string.
func FilesInFolderExt(path string, ext string) (files []string, err error) {
	logger.Verbose(fmt.Sprintf("[*] Searching for files in %s", path))

	useFilter := false
	if ext != "" {
		// normalize since Windows file extensions are not case sensitive
		ext = strings.ToLower(ext)
		useFilter = true
	}

	// Walk over every file in a given folder
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if useFilter {
			// Only append if it's not a directory and if it has the file extension specified
			if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ext {
				// Append to slice
				files = append(files, path)
			}
		} else {
			// Only append if it's not a directory
			if !info.IsDir() {
				// Append to slice
				files = append(files, path)
			}
		}

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
