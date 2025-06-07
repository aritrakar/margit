package main

type Blob struct {
	Data []byte
	Hash []byte
}

// type EntryType byte

// const (
// 	TypeBlob EntryType = iota
// 	TypeTree
// )

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
	ParentHash []byte
	Message    string
	Timestamp  string // for now, keep it string
	Hash       []byte
}
