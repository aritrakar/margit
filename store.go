package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func saveObject(obj interface{}) ([]byte, error) {
	// Serialize the object (e.g., Tree or Blob)
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	hashStr := fmt.Sprintf("%x", hash[:])

	// Check if the object hash already exists
	filePath := filepath.Join(OBJECT_DIR, hashStr)
	if _, err := os.Stat(filePath); err == nil {
		// Object already exists, return the existing hash
		return hash[:], nil
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return nil, err
	}

	return hash[:], nil
}

func loadObject(hash []byte, out interface{}) error {
	hashStr := fmt.Sprintf("%x", hash)
	filePath := filepath.Join(OBJECT_DIR, hashStr)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}
