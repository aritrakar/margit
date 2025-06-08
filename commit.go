package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NOTE: No support for branches yet so we use a fixed path for main branch.

const headPath = ".margit/HEAD"

func getCurrentRefPath() (string, error) {
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	// Trim whitespace and check if it starts with "ref: "
	line := strings.TrimSpace(string(data))
	if strings.HasPrefix(line, "ref: ") {
		return filepath.Join(".margit", strings.TrimPrefix(line, "ref: ")), nil
	}
	return "", fmt.Errorf("HEAD is not a ref")
}

func runInit() error {
	err := os.MkdirAll(".margit/objects", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create object directory: %w", err)
	}

	// Create refs directory
	err = os.MkdirAll(".margit/refs/heads", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create refs directory: %w", err)
	}

	// Create refs/heads/main file
	err = os.WriteFile(".margit/refs/heads/main", []byte{}, 0644)
	if err != nil {
		return fmt.Errorf("failed to create main branch reference: %w", err)
	}

	// HEAD file points to the main branch
	err = os.WriteFile(headPath, []byte("ref: refs/heads/main"), 0644)

	fmt.Println("Initialized empty Margit repository in .margit")
	return nil
}

func runCommit(message string) error {
	// Step 1: Build the tree from current directory
	// TODO: This needs to change to use a specific directory?
	tree, err := BuildTree("./test")
	if err != nil {
		return err
	}

	// Step 2: Load parent hash if available
	var parentHash []byte
	refPath, err := getCurrentRefPath()
	if err != nil {
		return fmt.Errorf("failed to get current ref path: %w", err)
	}

	if data, err := os.ReadFile(refPath); err == nil {
		parentHash = bytes.TrimSpace(data)
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
	err = os.WriteFile(refPath, commit.Hash, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Committed: %x\n", commit.Hash)
	return nil
}
