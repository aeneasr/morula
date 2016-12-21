package src

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)


// GetAfterAll returns the "after-all" section of the configuration file
func GetAfterAll() (result []string) {
	result = getConfigSectionAsSlice("after-all")
	for _, directory := range result {
		if !DirectoryExists(directory) {
			fmt.Printf("The config file specifies to run subproject %s after all others,\nbut such a subproject does not exist.", directory)
			os.Exit(1)
		}
		if !IsDirectory(directory) {
			fmt.Printf("The config file specifies to run subproject %s after all others,\nbut this path is not a directory.", directory)
			os.Exit(1)
		}
	}
	return
}

// GetAlways returns the "always" part of the configuration file
func GetAlways() (result []string) {
	result = getConfigSectionAsSlice("always")
	for _, directory := range result {
		if !DirectoryExists(directory) {
			fmt.Printf("The config file specifies to always run subproject %s,\nbut such a subproject does not exist.", directory)
			os.Exit(1)
		}
		if !IsDirectory(directory) {
			fmt.Printf("The config file specifies to always run subproject %s,\nbut this path is not a directory.", directory)
			os.Exit(1)
		}
	}
	return
}

// GetBeforeAll returns the "before-all" part of the configuration file
func GetBeforeAll() (result []string) {
	result = getConfigSectionAsSlice("before-all")
	for _, directory := range result {
		if !DirectoryExists(directory) {
			fmt.Printf("The config file specifies to run subproject %s before all others,\nbut such a subproject does not exist.", directory)
			os.Exit(1)
		}
		if !IsDirectory(directory) {
			fmt.Printf("The config file specifies to run subproject %s before all others,\nbut this path is not a directory.", directory)
			os.Exit(1)
		}
	}
	return
}


// GetNever returns the "never" part of the configuration file
func GetNever() (result []string) {
	result = getConfigSectionAsSlice("never")
	for _, directory := range result {
		if !DirectoryExists(directory) {
			fmt.Printf("The config file specifies to never run subproject %s,\nbut such a subproject does not exist.", directory)
			os.Exit(1)
		}
		if !IsDirectory(directory) {
			fmt.Printf("The config file specifies to never run subproject %s,\nbut this path is not a directory.", directory)
			os.Exit(1)
		}
	}
	return
}


func getConfigSectionAsSlice(key string) (result []string) {
	afterAllString := viper.GetString(key)
	if afterAllString != "" {
		result = []string{afterAllString}
	} else {
		result = viper.GetStringSlice(key)
	}
	return
}
