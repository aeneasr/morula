package src

import "os"


// IsDirectory returns whether the directory with the given name
// exists and is a directory.
func IsDirectory(dirName string) bool {
	stats, err := os.Stat(dirName)
	Check(err)
	return stats.IsDir()
}
