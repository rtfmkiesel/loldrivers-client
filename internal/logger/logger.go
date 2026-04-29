package logger

import (
	"fmt"
	"os"
	"strings"
)

var PrintDebug bool

func Stdout(s string, args ...any)  { fmt.Fprintf(os.Stdout, s, args...) }
func Stderr(s string, args ...any)  { fmt.Fprintf(os.Stderr, s, args...) }
func Debug(s string, args ...any)   { log("DBG", s, args...) }
func Info(s string, args ...any)    { log("INF", s, args...) }
func Warning(s string, args ...any) { log("WAR", s, args...) }
func Errorf(s string, args ...any)  { log("ERR", s, args...) }
func Error(err error)               { log("ERR", "%s", err.Error()) }
func Fatalf(s string, args ...any)  { log("FTL", s, args...); os.Exit(1) }
func Fatal(err error)               { log("FTL", "%s", err.Error()); os.Exit(1) }

func log(level, s string, args ...any) {
	if level == "DBG" && !PrintDebug {
		return
	}

	msg := fmt.Sprintf(s, args...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	levelColor := "\033[37m"
	switch level {
	case "INF":
		levelColor = "\033[32m" // green
	case "WAR":
		levelColor = "\033[33m" // yellow
	case "ERR", "FTL":
		levelColor = "\033[31m" // red
	}

	fmt.Fprintf(os.Stdout, "%s[%s]\033[0m %s", levelColor, level, msg)
}
