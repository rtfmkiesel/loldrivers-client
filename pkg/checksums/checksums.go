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

// Struct for the return values of Calc()
type Sums struct {
	Filename string
	MD5      string
	SHA1     string
	SHA256   string
}

// Calc() will return the MD5, SHA1 and SHA256 checksum of a given file as a Sums struct
func Calc(filepath string) (sums Sums, err error) {
	// Calculate the MD5
	sums.MD5, err = MD5(filepath)
	if err != nil {
		return sums, err
	}

	// Calculate the SHA1
	sums.SHA1, err = SHA1(filepath)
	if err != nil {
		return sums, err
	}

	// Calculate the SHA256
	sums.SHA256, err = SHA256(filepath)
	if err != nil {
		return sums, err
	}

	// Set the filename
	sums.Filename = filepath

	return sums, nil
}

// MD5() will return the MD5 checksum of the given file
func MD5(filepath string) (string, error) {
	// Check if the file exists
	if !filesystem.FileExists(filepath) {
		return "", fmt.Errorf("file '%s' does not exist", filepath)
	}

	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		errormsg := fmt.Sprintf("%s", err)
		// Ignore open errors "cannot access the file", "file cannot be accessed", "Access denied"
		if strings.Contains(strings.ToLower(errormsg), "access") {
			return "", nil
		}

		// Ignore "file does not exist" error because files could have been removed in the meantime
		// os.IsExist does not work
		if strings.Contains(strings.ToLower(errormsg), "does not exist") {
			return "", nil
		}

		return "", err
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

// SHA1 will return the SHA1 checksum of the given file
func SHA1(filepath string) (string, error) {
	// Check if the file exists
	if !filesystem.FileExists(filepath) {
		return "", fmt.Errorf("file '%s' does not exist", filepath)
	}

	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		errormsg := fmt.Sprintf("%s", err)
		// Ignore open errors "cannot access the file", "file cannot be accessed", "Access denied"
		if strings.Contains(strings.ToLower(errormsg), "access") {
			return "", nil
		}

		// Ignore "file does not exist" error because files could have been removed in the meantime
		// os.IsExist does not work
		if strings.Contains(strings.ToLower(errormsg), "does not exist") {
			return "", nil
		}

		return "", err
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

// SHA256 will return the SHA256 checksum of the given file
func SHA256(filepath string) (string, error) {
	// Check if the file exists
	if !filesystem.FileExists(filepath) {
		return "", fmt.Errorf("file '%s' does not exist", filepath)
	}

	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		errormsg := fmt.Sprintf("%s", err)
		// Ignore open errors "cannot access the file", "file cannot be accessed", "Access denied"
		if strings.Contains(strings.ToLower(errormsg), "access") {
			return "", nil
		}

		// Ignore "file does not exist" error because files could have been removed in the meantime
		// os.IsExist does not work
		if strings.Contains(strings.ToLower(errormsg), "does not exist") {
			return "", nil
		}

		return "", err
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

// CalcRunner() is used as a go func for calculating file checksums from a job list of filenames
func CalcRunner(wg *sync.WaitGroup, chanJobs <-chan string, chanResults chan<- Sums) {
	defer wg.Done()

	// For each job
	for job := range chanJobs {
		// Calculate checksums
		sums, err := Calc(job)
		if err != nil {
			logger.Catch(err)
			continue
		}

		// Add to the output channel
		chanResults <- sums
	}
}

// CompareRunner() is used as a go func for comparing checksums
func CompareRunner(wg *sync.WaitGroup, chanJobs <-chan Sums, chanResults chan<- logger.Result, checksums loldrivers.DriverHashes) {
	defer wg.Done()

	// For each job
	for job := range chanJobs {
		// Check if the MD5 in in the driver slice
		if contains(checksums.MD5Sums, job.MD5) {
			// Send result to the output channel
			chanResults <- logger.Result{
				Filename: job.Filename,
				Checksum: job.MD5,
			}
			continue
		}
		// Check if the SHA1 in in the driver slice
		if contains(checksums.SHA1Sums, job.SHA1) {
			// Send result to the output channel
			chanResults <- logger.Result{
				Filename: job.Filename,
				Checksum: job.SHA1,
			}
			continue
		}
		// Check if the SHA256 in in the driver slice
		if contains(checksums.SHA256Sums, job.SHA256) {
			// Send result to the output channel
			chanResults <- logger.Result{
				Filename: job.Filename,
				Checksum: job.SHA256,
			}
			continue
		}
	}
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
