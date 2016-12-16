package src

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)


// GetAfterAll returns the "after-all" section of the configuration file
func GetAfterAll() (result string) {
	result = viper.GetString("after-all")
	if result != "" && !DirectoryExists(result) {
		fmt.Printf("The config file specifies to run subproject %s after all others,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !IsDirectory(result) {
		fmt.Printf("The config file specifies to run subproject %s after all others,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}

// GetAlways returns the "always" part of the configuration file
func GetAlways() (result string) {
	result = viper.GetString("always")
	if result != "" && !DirectoryExists(result) {
		fmt.Printf("The config file specifies to always run subproject %s,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !IsDirectory(result) {
		fmt.Printf("The config file specifies to always run subproject %s,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}

// GetBeforeAll returns the "before-all" part of the configuration file
func GetBeforeAll() (result string) {
	result = viper.GetString("before-all")
	if result != "" && !DirectoryExists(result) {
		fmt.Printf("The config file specifies to run subproject %s before all others,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !IsDirectory(result) {
		fmt.Printf("The config file specifies to run subproject %s before all others,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}


// GetNever returns the "never" part of the configuration file
func GetNever() (result string) {
	result = viper.GetString("never")
	if result != "" && !DirectoryExists(result) {
		fmt.Printf("The config file specifies to never run subproject %s,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !IsDirectory(result) {
		fmt.Printf("The config file specifies to never run subproject %s,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}
