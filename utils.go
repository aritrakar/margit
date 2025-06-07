package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

const (
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func ensureObjectDir() error {
	return os.MkdirAll(objectDir, os.ModePerm)
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

func printTreeFromHash(hash []byte, indent string) {
	var tree Tree
	err := loadObject(hash, &tree)
	if err != nil {
		fmt.Printf("%s(error loading tree: %v)\n", indent, err)
		return
	}

	for _, entry := range tree.Entries {
		fmt.Printf("%s- %s [%s] %x\n", indent, entry.Name, typeToString(entry.Type), entry.Hash)
		if entry.Type == EntryTree {
			printTreeFromHash(entry.Hash, indent+"  ")
		}
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
