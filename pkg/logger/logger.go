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

func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "[!] ERROR: %s\n", err)
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
