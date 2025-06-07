package main

import (
	"os"
	"path/filepath"
	"sort"
)

func BuildTree(path string) (*Tree, error) {
	var entries []TreeEntry

	// Read the directory contents at the given path
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Iterate through files and directories in the given path
	// and create TreeEntry objects for each.
	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())

		if file.IsDir() {
			subTree, err := BuildTree(fullPath)
			if err != nil {
				return nil, err
			}
			entries = append(entries, TreeEntry{
				Name: file.Name(),
				Type: EntryTree,
				Hash: subTree.Hash,
			})
		} else {
			blob, err := createBlob(fullPath)
			if err != nil {
				return nil, err
			}
			entries = append(entries, TreeEntry{
				Name: file.Name(),
				Type: EntryBlob,
				Hash: blob.Hash,
			})
		}
	}

	// TODO: This might be unnecessary.
	// Sort entries by name to ensure deterministic hashing
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	// Serialize entries to JSON
	// serialized, err := json.Marshal(entries)
	// if err != nil {
	// 	return nil, err
	// }
	// treeHash := sha256.Sum256(serialized)
	// return &Tree{Entries: entries, Hash: treeHash[:]}, nil

	tree := &Tree{Entries: entries}
	hash, err := saveObject(tree)
	if err != nil {
		return nil, err
	}

	tree.Hash = hash
	return tree, nil
}
