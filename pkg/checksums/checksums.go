package checksums

import (
	"crypto/md5"  //#nosec G501
	"crypto/sha1" //#nosec G505
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// Calculate the SHA2 / SHA256 of a file at filePath
//
// Returns the hexadecimal string of the checksum
func calcSHA256(filePath string) (string, error) {
	return calcChecksum(filePath, sha256.New())
}

// Calculate the SHA1 / SHA128 of a file at filePath
//
// Returns the hexadecimal string of the checksum
func calcSHA1(filePath string) (string, error) {
	return calcChecksum(filePath, sha1.New()) //#nosec G401

}

// Calculate the MD5 of a file at filePath
//
// Returns the hexadecimal string of the checksum
func calcMD5(filePath string) (string, error) {
	return calcChecksum(filePath, md5.New()) //#nosec G401
}

func calcChecksum(filePath string, h hash.Hash) (string, error) {
	file, err := os.Open(filePath) //#nosec G304
	if err != nil {
		return "", err
	}
	defer file.Close() //nolint:errcheck

	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
