package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	Verbose      bool = false // Silences messages to stderr
	styleInfo         = "[" + color.BlueString("INF") + "]"
	styleWarning      = "[" + color.YellowString("WAR") + "]"
	styleError        = "[" + color.RedString("ERR") + "]"
	styleFatal        = "[" + color.RedString("FTL") + "]"
)

// Print a plain message to stdout
func PlainStdout(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	printStdout(msg)
}

// Print a plain message to stderr
func PlainStderr(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	printStderr(msg)
}

// Print a info message to stderr
func Info(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	msg = fmt.Sprintf("%s %s", styleInfo, msg)
	printStderr(msg)
}

// Print a warning message to stderr
func Warning(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	msg = fmt.Sprintf("%s %s", styleWarning, msg)
	printStderr(msg)
}

// Print an error to stderr
func Error(err error) {
	msg := fmt.Sprintf("%s %s", styleError, err)
	printStderr(msg)
}

// Print an error to stderr and quit
func Fatal(err error) {
	msg := fmt.Sprintf("%s %s", styleFatal, err)
	printStderr(msg)
	os.Exit(1)
}

// Print a message with a custom label to stdout
func Custom(labelStr string, colorAttrib color.Attribute, s string, args ...interface{}) {
	label := fmt.Sprintf("[%s]", color.New(colorAttrib).Sprint(labelStr))

	msg := fmt.Sprintf(s, args...)
	msg = fmt.Sprintf("%s %s", label, msg)
	msg = addNewLineIfNot(msg)
	printStdout(msg)
}

func addNewLineIfNot(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s
	} else {
		return s + "\n"
	}
}

func printStdout(s string) {
	fmt.Fprint(color.Output, addNewLineIfNot(s))
}

func printStderr(s string) {
	if Verbose {
		fmt.Fprint(color.Error, addNewLineIfNot(s))
	}
}
