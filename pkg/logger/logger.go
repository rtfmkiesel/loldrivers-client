package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	ShowDebugOutput = false // If set to true, logger.Debug() messages will get printed
	styleDebug      = "[" + color.BlueString("DBG") + "]"
	styleError      = "[" + color.RedString("ERR") + "]"
	styleFatal      = "[" + color.RedString("FTL") + "]"
)

// Prints a message with a custom label to stdout
func Custom(labelStr string, colorAttrib color.Attribute, s string, args ...any) {
	label := fmt.Sprintf("[%s]", color.New(colorAttrib).Sprint(labelStr))

	msg := fmt.Sprintf(s, args...)
	msg = fmt.Sprintf("%s %s", label, msg)
	fmt.Fprint(color.Output, formatMsg(msg)) //nolint:errcheck
}

// Prints a format string to stderr
func Stderr(s string, args ...any) {
	msg := fmt.Sprintf(s, args...)
	fmt.Fprint(color.Error, formatMsg(msg)) //nolint:errcheck
}

// Prints a format string to stdout
func Stdout(s string, args ...any) {
	msg := fmt.Sprintf(s, args...)
	fmt.Fprint(color.Output, formatMsg(msg)) //nolint:errcheck
}

// Prints a debug/info format string to stderr
func Debug(s string, args ...any) {
	if !ShowDebugOutput {
		return
	}

	msg := fmt.Sprintf(s, args...)
	msg = fmt.Sprintf("%s %s", styleDebug, msg)

	fmt.Fprint(color.Error, formatMsg(msg)) //nolint:errcheck
}

// Prints an error to stderr
func Error(err error) {
	msg := fmt.Sprintf("%s %s", styleError, err)
	fmt.Fprint(color.Error, formatMsg(msg)) //nolint:errcheck
}

// Prints an error format string to stderr
func Errorf(s string, args ...any) {
	Error(fmt.Errorf(s, args...))
}

// Prints an fatal error to stderr and quits
func Fatal(err error) {
	msg := fmt.Sprintf("%s %s", styleFatal, err)
	fmt.Fprint(color.Error, formatMsg(msg)) //nolint:errcheck
	os.Exit(1)
}

func formatMsg(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s
	} else {
		return s + "\n"
	}
}
