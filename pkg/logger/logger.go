// Package logger handles errors as well as the output handling
package logger

import (
	"fmt"
	"os"
	"strings"
)

var (
	BeSilent bool = false
)

func Log(msg string) {
	if !BeSilent {
		if strings.HasSuffix(msg, "\n") {
			fmt.Fprint(os.Stdout, msg)
		} else {
			fmt.Fprint(os.Stdout, msg+"\n")
		}
	}
}

func Logf(msg string, args ...interface{}) {
	if !BeSilent {
		msg = fmt.Sprintf(msg, args...)
		if strings.HasSuffix(msg, "\n") {
			fmt.Fprint(os.Stdout, msg)
		} else {
			fmt.Fprint(os.Stdout, msg+"\n")
		}
	}
}

func Error(err error) {
	fmt.Fprintf(os.Stderr, "[!] ERROR: %s\n", err)
}

func Errorf(msg string, args ...interface{}) {
	msg = fmt.Sprintf("[!] ERROR: "+msg, args...)
	if strings.HasSuffix(msg, "\n") {
		fmt.Fprint(os.Stderr, msg)
	} else {
		fmt.Fprint(os.Stderr, msg+"\n")
	}
}

func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "[!] ERROR: %s\n", err)
}

func Fatalf(msg string, args ...interface{}) {
	msg = fmt.Sprintf("[!] ERROR: "+msg, args...)
	if strings.HasSuffix(msg, "\n") {
		fmt.Fprint(os.Stderr, msg)
	} else {
		fmt.Fprint(os.Stderr, msg+"\n")
	}
	os.Exit(1)
}

func Banner() {
	if !BeSilent {
		fmt.Printf(`
    ┌─────────────────────────────────────┐
    │          LOLDrivers-client          │
    │      https://www.loldrivers.io      │
    │                                     │
    │    https://github.com/rtfmkiesel    │ 
    └─────────────────────────────────────┘

`)
	}
}
