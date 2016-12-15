package src

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

type Runner struct {
	C aurora.Aurora
	Command string
}


func NewRunner(C aurora.Aurora, commands []string) *Runner {
	return &Runner{C, strings.Join(commands, " ")}
}


func (runner *Runner) RunInSubproject(subprojectName string) (err error) {

	// determine directory to run the command in
	dir := runner.getDirectoryToRunIn(subprojectName)

	// run the command
	fmt.Printf("running %s in subproject %s ...\n\n", runner.C.Bold(runner.C.Cyan(runner.Command)), runner.C.Bold(runner.C.Cyan(subprojectName)))
	cmd := runner.createCommand()
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	Check(err) // this error should always be nil, since we call the command shell which always exists
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("subproject %s has issues\n", subprojectName)
		return err
	}

	fmt.Print("\n\n")
	return
}


func (runner *Runner) createCommand() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", runner.Command)
	} else {
		return exec.Command("bash", "-c", runner.Command)
	}
}


// determine directory to run the command in
func (runner *Runner) getDirectoryToRunIn(subprojectName string) string {
	cwd, err := os.Getwd()
	Check(err)
	return path.Join(cwd, subprojectName)
}
