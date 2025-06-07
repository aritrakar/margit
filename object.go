package main

import "time"

type Blob struct {
	Data []byte
	Hash []byte
}

const (
	EntryBlob = 0
	EntryTree = 1
)

type TreeEntry struct {
	Name string
	Type uint8 // 0 for blob, 1 for tree
	Hash []byte
}

type Tree struct {
	Entries []TreeEntry
	Hash    []byte
}

type Commit struct {
	TreeHash   []byte
	ParentHash []byte // nil for initial commit
	Message    string
	Timestamp  time.Time
	Hash       []byte
}
