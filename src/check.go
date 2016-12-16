package src

import (
	"fmt"
	"os"
)

// Check allows convenient checking of critical errors.
// When an error is given, it prints it and ends the application.
func Check(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
		os.Exit(1)
	}
}
