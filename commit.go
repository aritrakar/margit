package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func runCommit(message string) error {
	err := ensureObjectDir()
	if err != nil {
		return err
	}

	// Step 1: Build the tree from current directory
	tree, err := BuildTree("./test")
	if err != nil {
		return err
	}

	// Step 2: Load parent hash if available
	var parentHash []byte
	headPath := ".margit/HEAD"
	if data, err := os.ReadFile(headPath); err == nil {
		parentHash = data
	}

	// Step 3: Create and save commit object
	commit := &Commit{
		TreeHash:   tree.Hash,
		ParentHash: parentHash,
		Message:    message,
		Timestamp:  time.Now(),
	}

	// Serialize and hash
	raw, err := json.Marshal(commit)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(raw)
	commit.Hash = hash[:]

	// Save commit object
	err = os.WriteFile(filepath.Join(objectDir, fmt.Sprintf("%x", commit.Hash)), raw, 0644)
	if err != nil {
		return err
	}

	// Step 4: Update HEAD
	err = os.WriteFile(headPath, commit.Hash, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Committed: %x\n", commit.Hash)
	return nil
}
