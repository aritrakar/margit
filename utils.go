package main

import (
	"fmt"
	"os"
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
