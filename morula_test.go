package main

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kr/pretty"
	"github.com/termie/go-shutil"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func FeatureContext(s *godog.Suite) {

	// the temp dir in which the test repos live
	var testRoot string

	// the output of the last command run
	var output string

	// the error of the last run operation
	var err error

	s.BeforeScenario(func(interface{}) {
		fmt.Println("BEFORE ALL")
	})

	s.AfterScenario(func(interface{}, error) {
		fmt.Println("AFTER ALL")
	})

	s.Step(`^a project with the subprojects:$`, func(projectData *gherkin.DataTable) error {
		testRoot = createTempDir()
		for _, project := range projectData.Rows[1:] {
			createTestProject(project.Cells[1].Value, project.Cells[0].Value, testRoot)
		}
		return nil
	})

	s.Step(`^it fails with an error code and the message:$`, func(expectedText *gherkin.DocString) error {
		if err == nil {
			return errors.New("Expected error, but test passed")
		}
		if !strings.Contains(output, expectedText.Content) {
			return errors.New(fmt.Sprintf("Expected to see \"%s\" in \"%s\"", expectedText.Content, output))
		}
		return nil
	})

	s.Step(`^it runs the tests:$`, func(tests *gherkin.DataTable) (result error) {

		// determine the names of the projects we expect to be tested
		var expectedProjectNames []string
		for _, project := range tests.Rows {
			expectedProjectNames = append(expectedProjectNames, project.Cells[0].Value)
		}

		// determine the projects morula has actually tested
		re := regexp.MustCompile(`testing subproject (.*) \.\.\.`)
		matches := re.FindAllStringSubmatch(output, -1)
		var actualProjectNames = make([]string, len(matches))
		for i, match := range matches {
			actualProjectNames[i] = match[1]
		}

		diffs := pretty.Diff(expectedProjectNames, actualProjectNames)
		if len(diffs) > 0 {
			result = fmt.Errorf("Didn't run the expected projects: %s", diffs[0])
		}
		return
	})

	s.Step(`^running "([^"]*)"$`, func(command string) (result error) {
		output, result = run(command, testRoot)
		fmt.Println(output)
		return
	})

	s.Step(`^trying to run "([^"]*)"$`, func(command string) (result error) {
		output, err = run(command, testRoot)
		fmt.Println(output)
		if err == nil {
			result = errors.New("Expected failure, but command ran without errors")
		}
		return
	})

}

func createTempDir() (path string) {
	dir, err := ioutil.TempDir("", "morula")
	check(err)
	return dir
}

// Creates a subproject of the given type,
// with the given name,
// in the given test workspace
func createTestProject(template string, name string, testRoot string) {
	check(shutil.CopyTree(
		filepath.Join("features", "examples", template),
		filepath.Join(testRoot, name),
		nil))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Runs the given command, returns its output
func run(commandText string, dir string) (output string, err error) {
	args := strings.Split(commandText, " ")
	args = append(args, "--color=false")
	command := exec.Command(args[0], args[1:]...)
	command.Dir = dir
	outputArray, err := command.CombinedOutput()
	output = string(outputArray)
	return
}

// Returns the project names in the given human-friendly project name list
func splitProjectNames(projectNames string) (result []string) {
	for _, projectName := range strings.Split(projectNames, "and") {
		result = append(result, strings.Trim(projectName, " \""))
	}
	return
}
