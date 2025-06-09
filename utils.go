package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func ensureObjectDir() error {
	return os.MkdirAll(OBJECT_DIR, os.ModePerm)
}

func typeToString(t uint8) string {
	switch t {
	case EntryBlob:
		return "blob"
	case EntryTree:
		return "tree"
	default:
		return "unknown"
	}
}

func test(path string) {
	tree, err := BuildTree(path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Merkle root of directory:", hex.EncodeToString(tree.Hash))
	fmt.Println("Tree entries:")
	printTreeFromHash(tree.Hash, "")
}

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

// func flattenWorkingDir(root string, prefix string, result map[string][32]byte) error {
// 	fmt.Println("Flattening working directory:", root)

// 	files, err := os.ReadDir(root)
// 	if err != nil {
// 		return err
// 	}

// 	for _, file := range files {
// 		name := file.Name()
// 		diskPath := filepath.Join(root, name)
// 		logicalPath := filepath.Join(prefix, name)

// 		if file.IsDir() {
// 			if err := flattenWorkingDir(diskPath, logicalPath, result); err != nil {
// 				return err
// 			}
// 		} else {
// 			data, err := os.ReadFile(diskPath)
// 			if err != nil {
// 				return err
// 			}
// 			result[logicalPath] = sha256.Sum256(data)
// 		}
// 	}
// 	return nil
// }
