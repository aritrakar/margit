package main

import "os"

/*
Creates a Blob object from a file at the given path.
*/

func createBlob(path string) (*Blob, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	blob := &Blob{Data: data}

	// Save the blob object to the object store and get its hash
	hash, err := saveObject(blob)
	if err != nil {
		return nil, err
	}

	blob.Hash = hash
	return blob, nil
}
