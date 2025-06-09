package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// NOTE: No support for branches yet so we use a fixed path for main branch.

const headPath = ".margit/HEAD"

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

	// TODO: Support branches
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
	tree, err := BuildTree(WORKING_DIR)
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
	err = os.WriteFile(filepath.Join(OBJECT_DIR, fmt.Sprintf("%x", commit.Hash)), raw, 0644)
	if err != nil {
		return err
	}

	// Step 4: Update HEAD
	err = os.WriteFile(refPath, commit.Hash, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Committed: %x\n", commit.Hash)

	// Debug. Print all hashes in the current commit
	fmt.Println("Current commit tree:")
	treeMap := make(map[string][32]byte)
	if err := flattenTree(tree, "", treeMap); err != nil {
		fmt.Println("Error flattening tree:", err)
	}

	for path, hash := range treeMap {
		fmt.Printf("%s -> %x\n", path, hash)
	}

	return nil
}

func getLatestCommit() (*Commit, error) {
	refPath, err := getCurrentRefPath()
	if err != nil {
		// No HEAD ref available, no commits yet.
		return nil, nil
	}

	commitHash, err := os.ReadFile(refPath)
	if err != nil {
		// No commit file found.
		return nil, fmt.Errorf("no commits yet")
	}

	var commit Commit
	if err := loadObject(bytes.TrimSpace(commitHash), &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

func runStatus() error {
	// Get the latest commit
	commit, err := getLatestCommit()
	if err != nil {
		fmt.Println("Error getting latest commit:", err)
		return err
	}

	// Load the tree object for the current commit (i.e. without current changes)
	var tree Tree
	if err := loadObject(commit.TreeHash, &tree); err != nil {
		return err
	}

	// Flatten the tree into a map of paths to hashes
	// This will help us compare with the working directory
	committed := make(map[string][32]byte)
	if err := flattenTree(&tree, "", committed); err != nil {
		return err
	}

	// // Flatten the working directory into a map of paths to hashes
	// working := make(map[string][32]byte)
	// if err := flattenWorkingDir(WORKING_DIR, "", working); err != nil {
	// 	return err
	// }

	// TODO: Find a way to make this more efficient.

	// Build the working directory tree
	workingTree, err := BuildTree(WORKING_DIR)
	if err != nil {
		fmt.Println("Error building working tree:", err)
		os.Exit(1)
	}

	// Flatten the working directory tree into a map of paths to hashes
	working := make(map[string][32]byte)
	if err := flattenTree(workingTree, "", working); err != nil {
		fmt.Println("Error flattening working tree:", err)
		os.Exit(1)
	}

	fmt.Println("On branch main")
	fmt.Println("Changes not staged for commit:")

	for path, hash := range working {
		committedHash, exists := committed[path]
		if !exists {
			fmt.Printf("  added: %s\n", path)
		} else if hash != committedHash {
			fmt.Printf("  modified: %s\n", path)
			// fmt.Printf("  working hash  : %x\n", hash)
			// fmt.Printf("  committed hash: %x\n", committedHash)
		}
	}

	for path := range committed {
		if _, exists := working[path]; !exists {
			fmt.Printf("  deleted: %s\n", path)
		}
	}

	return nil
}
