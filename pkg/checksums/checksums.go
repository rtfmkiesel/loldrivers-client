// Package checksums calculates the MD5, SHA1 and SHA256 checksum of files
package checksums

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"loldrivers-client/pkg/filesystem"
	"loldrivers-client/pkg/logger"
	"loldrivers-client/pkg/loldrivers"
)

// calcMD5() will return the MD5 checksum of the given file
func calcMD5(filepath string) (string, error) {
	// Check if the file exists
	if !filesystem.FileExists(filepath) {
		return "", fmt.Errorf("file '%s' does not exist", filepath)
	}

	// Open the file
	file, err := os.Open(filepath)
	// Check if the file can be accessed
	if fileAccessErr(err) {
		// No, skip file
		return "", nil
	}
	defer file.Close()

	// Create a new MD5
	hash := md5.New()

	// Copy the file data into the hash
	if _, err := io.Copy(hash, file); err != nil {
		errormsg := fmt.Sprintf("%s", err)
		// Ignore read error "another process has locked a portion of the file"
		if strings.Contains(strings.ToLower(errormsg), "has locked") {
			return "", nil
		}

		return "", err
	}

	// Get the checksum
	checksum := hash.Sum(nil)

	// Convert the checksum to a hex string
	return hex.EncodeToString(checksum), nil
}

// calcSHA1 will return the SHA1 checksum of the given file
func calcSHA1(filepath string) (string, error) {
	// Check if the file exists
	if !filesystem.FileExists(filepath) {
		return "", fmt.Errorf("file '%s' does not exist", filepath)
	}

	// Open the file
	file, err := os.Open(filepath)
	// Check if the file can be accessed
	if fileAccessErr(err) {
		// No, skip file
		return "", nil
	}
	defer file.Close()

	// Create a new SHA1
	hash := sha1.New()

	// Copy the file data into the hash
	if _, err := io.Copy(hash, file); err != nil {
		errormsg := fmt.Sprintf("%s", err)
		// Ignore read error "another process has locked a portion of the file"
		if strings.Contains(strings.ToLower(errormsg), "has locked") {
			return "", nil
		}

		return "", err
	}

	// Get the checksum
	checksum := hash.Sum(nil)

	// Convert the checksum to a hex string
	return hex.EncodeToString(checksum), nil
}

// calcSHA256 will return the SHA256 checksum of the given file
func calcSHA256(filepath string) (string, error) {
	// Check if the file exists
	if !filesystem.FileExists(filepath) {
		return "", fmt.Errorf("file '%s' does not exist", filepath)
	}

	// Open the file
	file, err := os.Open(filepath)
	// Check if the file can be accessed
	if fileAccessErr(err) {
		// No, skip file
		return "", nil
	}
	defer file.Close()

	// Create a new SHA256
	hash := sha256.New()

	// Copy the file data into the hash
	if _, err := io.Copy(hash, file); err != nil {
		errormsg := fmt.Sprintf("%s", err)
		// Ignore read error "another process has locked a portion of the file"
		if strings.Contains(strings.ToLower(errormsg), "has locked") {
			return "", nil
		}

		return "", err
	}

	// Get the checksum
	checksum := hash.Sum(nil)

	// Convert the checksum to a hex string
	return hex.EncodeToString(checksum), nil
}

// contains() will return true if a []string contains a string
func contains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}

	return false
}

// fileAccessErr() will handle ACL errors.
// It will return true if a file should be skipped (can't be read/opened)
func fileAccessErr(err error) bool {
	if err == nil {
		return false
	}

	errormsg := fmt.Sprintf("%s", err)
	// Ignore open errors "cannot access the file", "file cannot be accessed", "Access denied"
	if strings.Contains(strings.ToLower(errormsg), "access") {
		return true
	}

	// Ignore "file does not exist" error because files could have been removed in the meantime
	// os.IsExist does not work
	if strings.Contains(strings.ToLower(errormsg), "does not exist") {
		return true
	}

	return false
}

// Runner() is used as a go func for calculating and comparing file checksums
// from a job channel of filenames. If a calculated checksum matches a loaded checksum
// a result in the form of logger.Result will be sent to an output channel
func Runner(wg *sync.WaitGroup, chanJobs <-chan string, chanResults chan<- logger.Result, checksums loldrivers.DriverHashes) {
	defer wg.Done()

	// For each job
	for job := range chanJobs {
		// Calculate the MD5
		MD5, err := calcMD5(job)
		if err != nil {
			logger.Catch(err)
			continue
		}
		// Check if the MD5 in in the driver slice
		if contains(checksums.MD5Sums, MD5) {
			// Send result to the output channel
			chanResults <- logger.Result{
				Filename: job,
				Checksum: MD5,
			}
			continue
		}

		// Calculate the SHA1
		SHA1, err := calcSHA1(job)
		if err != nil {
			logger.Catch(err)
			continue
		}
		// Check if the SHA1 in in the driver slice
		if contains(checksums.SHA1Sums, SHA1) {
			// Send result to the output channel
			chanResults <- logger.Result{
				Filename: job,
				Checksum: SHA1,
			}
			continue
		}

		// Calculate the SHA256
		SHA256, err := calcSHA256(job)
		if err != nil {
			logger.Catch(err)
			continue
		}
		// Check if the SHA256 in in the driver slice
		if contains(checksums.SHA256Sums, SHA256) {
			// Send result to the output channel
			chanResults <- logger.Result{
				Filename: job,
				Checksum: SHA256,
			}
			continue
		}
	}
}
