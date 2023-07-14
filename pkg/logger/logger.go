// Package logger handles errors as well as the output handling
package logger

import (
	"fmt"
	"os"
)

var (
	BeSilent bool = false
)

// Log() will print to the terminal
func Log(message string) {
	if !BeSilent {
		fmt.Printf("%s\n", message)
	}
}

// Catch() will handle errors
func Catch(err error) {
	fmt.Printf("ERROR: %s\n", err)
}

// CatchCrit() will handle critical errors
func CatchCrit(err error) {
	fmt.Printf("CRITICAL ERROR: %s\n", err)
	os.Exit(1)
}

// Banner() will print the banner
func Banner() {
	if !BeSilent {
		fmt.Printf(`
    ╒═══════════════════════════╕
    |     LOLDrivers-client     |
    |   https://loldrivers.io   |
    |                           |
    | by @rtfmkiesel/mkiesel.ch |
    ╘═══════════════════════════╛

`)
	}
}
