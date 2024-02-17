package checksums

import (
	"crypto/md5"  //#nosec G501
	"crypto/sha1" //#nosec G505
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Calculate the SHA2 / SHA256 of a file at filePath
//
// Returns the hexadecimal string of the checksum
func calcSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath) //#nosec G304
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}

// Calculate the SHA1 / SHA128 of a file at filePath
//
// Returns the hexadecimal string of the checksum
func calcSHA1(filePath string) (string, error) {
	file, err := os.Open(filePath) //#nosec G304
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New() //#nosec G401
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}

// Calculate the MD5 of a file at filePath
//
// Returns the hexadecimal string of the checksum
func calcMD5(filePath string) (string, error) {
	file, err := os.Open(filePath) //#nosec G304
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New() //#nosec G401
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}
