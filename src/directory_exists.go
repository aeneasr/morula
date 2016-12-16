package src

import (
	"os"
)

// DirectoryExists returns whether the given directory exists in the current working directory
func DirectoryExists(dirName string) bool {
	_, err := os.Stat(dirName)
	return !os.IsNotExist(err)
}
