package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/Originate/morula/src"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "morula",
	Short: "Optimizing task runner for monorepositories",
	Long: `Morula runs tasks in all subprojects of a monorepository.

The individual subprojects should be located in top-level folders.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.morula.yaml)")
	RootCmd.PersistentFlags().Bool("color", true, "Display output in color")
	RootCmd.PersistentFlags().String("always", "", "subproject to always run")
	RootCmd.PersistentFlags().String("never", "", "subproject to never run")
	RootCmd.PersistentFlags().String("after-all", "", "subproject to run after all others")
	RootCmd.PersistentFlags().String("before-all", "", "subproject to run before all others")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// viper.BindPFlag("color", RootCmd.PersistentFlags().Lookup("color"))
	viper.BindPFlag("always", RootCmd.PersistentFlags().Lookup("always"))
	viper.BindPFlag("never", RootCmd.PersistentFlags().Lookup("never"))
	viper.BindPFlag("after-all", RootCmd.PersistentFlags().Lookup("after-all"))
	viper.BindPFlag("before-all", RootCmd.PersistentFlags().Lookup("before-all"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("morula") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// HELPER FUNCTIONS

func getAurora(cmd *cobra.Command) aurora.Aurora {
	color, err := cmd.Flags().GetBool("color")
	src.Check(err)
	return aurora.NewAurora(color)
}
