package main

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"io/ioutil"
	// "os"
)

type morulaFeature struct {
	testRoot string
}

func (m *morulaFeature) aProjectWithSubprojects(arg1 string, arg2 string) error {
	dir, err := m.getTempDir()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	// os.Mkdir(dir)
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	m := &morulaFeature{}

	s.Step(`^a project with the subprojects "([^"]*)" and "([^"]*)"$`, m.aProjectWithSubprojects)

	s.Step(`^it runs the tests:$`, func(arg1 *gherkin.DataTable) error {
		return godog.ErrPending
	})

	s.Step(`^running "([^"]*)"$`, func(command string) error {
		fmt.Println("testing")
		return nil
	})

}

func (m *morulaFeature) getTempDir() (path string, err error) {
	return ioutil.TempDir("", "morula")
}

func createTestProject(name string) {
}
