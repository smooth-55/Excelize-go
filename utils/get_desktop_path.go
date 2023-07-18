package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPath() string {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to retrieve user's home directory:", err)
		return ""
	}

	// Open the home directory
	dir, err := os.Open(homeDir)
	if err != nil {
		fmt.Println("Failed to open user's home directory:", err)
		return ""
	}
	defer dir.Close()

	// Read the directory entries
	entries, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Failed to read directory entries:", err)
		return ""
	}

	// Iterate over the entries and filter directories
	var _dir string
	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name()
			dirPath := filepath.Join(homeDir, dirName)
			// fmt.Println(dirPath)
			if strings.Contains(dirName, "Desktop") {
				_dir = dirPath
			}
		}
	}
	return _dir
}
