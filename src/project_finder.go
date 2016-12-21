package src

import (
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// AllSubprojectNames returns the names of all subprojects,
// irrespective of whether they contain changes.
func AllSubprojectNames() (result []string) {
	beforeAll := GetBeforeAll()
	afterAll := GetAfterAll()
	never := GetNever()
	GetAlways() // to check for invalid "always" entry
	if len(beforeAll) > 0 {
		result = append(result, beforeAll...)
	}
	for _, entry := range filesystemEntries() {
		entryName := entry.Name()
		if !entry.IsDir() {
			continue
		}
		if entryName[0] == '.' {
			continue
		}
		if Contains(beforeAll, entryName) || Contains(afterAll, entryName) {
			continue
		}
		if Contains(never, entryName) {
			continue
		}
		result = append(result, entry.Name())
	}
	if len(afterAll) > 0 {
		result = append(result, afterAll...)
	}
	return result
}

// ChangedSubprojectNames returns the names of all subprojects
// that contain changes compared to the master branch.
func ChangedSubprojectNames() (result []string) {
	// Due to the lack of array methods like "uniq" in Golang,
	// this method iterates the filenames sorted alphabetically
	// and only appends the resulting project name to the result if the last element isn't it.
	beforeAll := GetBeforeAll()
	afterAll := GetAfterAll()
	always := GetAlways()
	never := GetNever()
	if len(beforeAll) > 0 {
		result = append(result, beforeAll...)
	}
	if len(always) > 0 {
		result = append(result, always...)
	}
	filePaths := runConsoleCommand("git", "diff", "--name-only", "master")
	sort.Strings(filePaths)
	for _, filePath := range filePaths {
		filePath = strings.Trim(filePath, " ")
		if len(filePath) == 0 {
			continue
		}
		projectName := strings.Split(filePath, "/")[0] // Git always returns "/" even on Windows
		if !isDir(projectName) {
			continue
		}
		if projectName[0] == '.' {
			continue
		}
		if Contains(afterAll, projectName) || Contains(beforeAll, projectName) {
			continue
		}
		if Contains(never, projectName) {
			continue
		}
		if len(result) > 0 && result[len(result)-1] == projectName {
			continue
		}
		result = append(result, projectName)
	}
	if len(afterAll) > 0 {
		result = append(result, afterAll...)
	}
	return result
}

func filesystemEntries() []os.FileInfo {
	entries, err := ioutil.ReadDir(".")
	Check(err)
	return entries
}

func isDir(fileName string) bool {
	file, err := os.Stat(fileName)
	return err == nil && file.IsDir()
}

// Runs the given command on the console
// Returns the result as an array of strings
func runConsoleCommand(name string, args ...string) []string {
	out, err := exec.Command(name, args...).Output()
	Check(err)
	return strings.Split(string(out), "\n")
}
