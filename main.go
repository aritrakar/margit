package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

const objectDir = ".margit/objects"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: margit <command> [options]")
		os.Exit(1)
	}

	// Having the error check here ensures the object directory
	// exists before any commands are run
	err := ensureObjectDir()
	if err != nil {
		fmt.Println("Error ensuring object directory:", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		initCmd.Parse(os.Args[2:])

		err := runInit()
		if err != nil {
			fmt.Println("Initialization failed:", err)
			os.Exit(1)
		}
		fmt.Println("Margit repository initialized successfully.")
	case "commit":
		commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
		message := commitCmd.String("m", "", "Commit message")

		commitCmd.Parse(os.Args[2:])
		if *message == "" {
			fmt.Println("Error: Commit message required with -m")
			os.Exit(1)
		}

		err := runCommit(*message)
		if err != nil {
			fmt.Println("Commit failed:", err)
			os.Exit(1)
		}
	case "log":
		logCmd := flag.NewFlagSet("log", flag.ExitOnError)
		// No options for log yet, but can be extended later
		logCmd.Parse(os.Args[2:])

		refPath, err := getCurrentRefPath()
		if err != nil {
			fmt.Printf("failed to get current ref path: %w", err)
			os.Exit(1)
		}

		data, err := os.ReadFile(refPath)
		if err != nil {
			fmt.Println("Error reading HEAD:", err)
			os.Exit(1)
		}

		fmt.Println("Current HEAD:", hex.EncodeToString(data))
		fmt.Println("Commit history:")

		// Print the commit history by using the parent hashes
		currentHash := data
		for len(currentHash) > 0 {
			var commit Commit
			err := loadObject(currentHash, &commit)
			if err != nil {
				fmt.Println("Error loading commit:", err)
				break
			}

			fmt.Printf("%scommit %x%s\n", Yellow, currentHash, Reset)

			// Format timestamp like: Sat Jun 7 18:58:32 2025 -0400
			fmt.Printf("Date:   %s\n", commit.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"))

			fmt.Println()
			fmt.Printf("    %s\n", commit.Message)
			fmt.Println()

			if len(commit.ParentHash) == 0 {
				break
			}
			currentHash = commit.ParentHash
		}

	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
