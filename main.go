package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

const objectDir = ".margit/objects"

func main() {
	if err := ensureObjectDir(); err != nil {
		fmt.Println("Error creating object directory:", err)
		os.Exit(1)
	}

	tree, err := BuildTree("./test/test2")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Merkle root of directory:", hex.EncodeToString(tree.Hash))
	fmt.Println("Tree entries:")
	printTreeFromHash(tree.Hash, "")
}
