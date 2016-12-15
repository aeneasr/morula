package src

import (
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type ProjectFinder struct {
	Always, Never, BeforeAll, AfterAll string
}


func (p *ProjectFinder) AllSubprojectNames() (result []string) {
	if p.BeforeAll != "" {
		result = append(result, p.BeforeAll)
	}
	for _, entry := range filesystemEntries() {
		entryName := entry.Name()
		if !entry.IsDir() { continue }
		if entryName[0] == '.' { continue }
		if entryName == p.BeforeAll || entryName == p.AfterAll { continue }
		if entryName == p.Never { continue }
		result = append(result, entry.Name())
	}
	if p.AfterAll != "" {
		result = append(result, p.AfterAll)
	}
	return result
}

func (p *ProjectFinder) ChangedSubprojectNames() (result []string) {
	// Due to the lack of array methods like "uniq" in Golang,
	// this method iterates the filenames sorted alphabetically
	// and only appends the resulting project name to the result if the last element isn't it.
	if p.BeforeAll != "" {
		result = append(result, p.BeforeAll)
	}
	if p.Always != "" {
		result = append(result, p.Always)
	}
	filePaths := runConsoleCommand("git", "diff", "--name-only", "master")
	sort.Strings(filePaths)
	for _, filePath := range filePaths {
		filePath = strings.Trim(filePath, " ")
		if len(filePath) == 0 { continue }
		projectName := strings.Split(filePath, "/")[0] // Git always returns "/" even on Windows
		if !isDir(projectName) { continue }
		if projectName[0] == '.' { continue }
		if projectName == p.AfterAll || projectName == p.BeforeAll { continue }
		if projectName == p.Never { continue }
		if len(result) > 0 && result[len(result)-1] == projectName { continue }
		result = append(result, projectName)
	}
	if p.AfterAll != "" {
		result = append(result, p.AfterAll)
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
