package src

import "os"


func IsDirectory(dirName string) bool {
	stats, err := os.Stat(dirName)
	Check(err)
	return stats.IsDir()
}
