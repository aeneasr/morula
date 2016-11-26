package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// setupCmd represents the test command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates an example configuration file with the default options",
	Long:  `Creates a configuration file in YML format, containing all possible options and their default values.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getAurora(cmd)
		createConfigFile()
		fmt.Printf("configuration file %s created", c.Bold(c.Cyan("morula.yml")))
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}

func createConfigFile() {
	check(ioutil.WriteFile("morula.yml", []byte(`before-all: ""
after-all: ""
always: ""
never: ""
color: true`), 0644))
}
