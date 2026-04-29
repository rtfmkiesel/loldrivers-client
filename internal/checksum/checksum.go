package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

type Hashfunc struct {
	Name      string
	Calculate func(string) (string, error)
}

var Algorithms = []Hashfunc{
	{"SHA1", SHA1},
	{"SHA256", SHA256},
	{"MD5", MD5},
}

func SHA256(file string) (string, error) { return calcChecksum(file, sha256.New()) }
func SHA1(file string) (string, error)   { return calcChecksum(file, sha1.New()) }
func MD5(file string) (string, error)    { return calcChecksum(file, md5.New()) }

func calcChecksum(file string, h hash.Hash) (string, error) {
	fh, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fh.Close() //nolint:errcheck

	if _, err := io.Copy(h, fh); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
