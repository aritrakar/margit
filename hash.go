package main

import (
	"crypto/sha256"
	"os"
)

/*
Reads a file from the given path, computes its SHA-256 hash,
and returns the file's data and its hash.
*/
func hashFile(path string) ([]byte, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	hash := sha256.Sum256(data)
	return data, hash[:], nil
}
