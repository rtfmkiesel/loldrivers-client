package options

import (
	"errors"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
)

var (
	// The three default Windows driver directories
	windowsDriverDirs = []string{"C:\\Windows\\System32\\drivers", "C:\\Windows\\System32\\DriverStore\\FileRepository", "C:\\WINDOWS\\inf"}
)

type Options struct {
	Mode              string    // Mode to load the drivers (online, local, internal)
	ModeLocalFilePath string    // Filepath to *.json if mode == local
	ScanDirectories   []string  // Directories to scan for drivers
	ScanSizeLimit     int       // File size limit in MB
	ScanWorkers       int       // Amount of scan (checksum calculation) workers
	ScanShowErrors    bool      // Display errors during checksum calculation (often file read errors)
	OutputMode        string    // How to print the results (default, json, grep)
	StartTime         time.Time // To track execution time
}

// Parse the command line options into an Options struct
func Parse() (*Options, error) {
	opt := &Options{
		StartTime: time.Now(),
	}

	var (
		flagDir      string
		flagGrepable bool
		flagJson     bool
	)

	flag.StringVarP(&opt.Mode, "mode", "m", "online", "Operating Mode {online, local, internal}")
	flag.StringVarP(&opt.ModeLocalFilePath, "driver-file", "f", "", "File path to 'drivers.json', when mode == local")
	flag.StringVarP(&flagDir, "scan-dir", "d", "", "Directory to scan for drivers (default: Windows driver folders)")
	flag.IntVarP(&opt.ScanSizeLimit, "scan-size", "l", 10, "Size limit for files to scan in MB")
	flag.IntVarP(&opt.ScanWorkers, "workers", "w", 20, "Number of checksum \"threads\" to spawn")
	flag.BoolVarP(&opt.ScanShowErrors, "surpress-errors", "s", false, "Do not show file read errors when calculating checksums")
	flag.BoolVarP(&flagGrepable, "grepable", "g", false, "Will only output found files for easy parsing")
	flag.BoolVarP(&flagJson, "json", "j", false, "Format output as JSON")

	flag.Parse()

	switch opt.Mode {
	case "online", "internal":
		// we good
	case "local":
		if opt.ModeLocalFilePath == "" {
			return nil, errors.New("-m/--mode 'local' requires '-f/--driver-file'")
		}
	default:
		return nil, errors.New("invalid mode")
	}

	// Only one output style
	if flagGrepable && flagJson {
		return nil, errors.New("only use '-g/--grepable' or '-j/--json', not both")
	}

	switch {
	case flagGrepable:
		opt.OutputMode = "grep"
		logger.ShowDebugOutput = false
	case flagJson:
		opt.OutputMode = "json"
		logger.ShowDebugOutput = false
	default:
		logger.ShowDebugOutput = true
	}

	if flagDir != "" {
		// User specified a custom folder to scan
		opt.ScanDirectories = []string{flagDir}
	} else {
		// User did not specify a path with '-d', use the default Windows directories
		opt.ScanDirectories = windowsDriverDirs
	}

	printBanner()

	return opt, nil
}

func printBanner() {
	logger.Stderr(`
  ╔─────────────────────────────────────╗
  │          LOLDrivers-client          │
  │      https://www.loldrivers.io      │
  │                                     │
  │    https://github.com/rtfmkiesel    │ 
  ╚─────────────────────────────────────╝

`)
}
