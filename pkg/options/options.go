package options

import (
	"errors"
	"time"

	"github.com/projectdiscovery/goflags"
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
func Parse() (opt *Options, err error) {
	opt = &Options{}
	opt.StartTime = time.Now()

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("A tool to scan your computer for known vulnerable and known malicious Windows drivers")

	flagSet.CreateGroup("Mode", "Operating Mode",
		flagSet.StringVarP(&opt.Mode, "mode", "m", "online", "Operating Mode {online, local, internal}"),
		flagSet.StringVarP(&opt.ModeLocalFilePath, "driver-file", "f", "", "File path to 'drivers.json', when mode == local"),
	)

	var flagDir string
	flagSet.CreateGroup("Scan", "Scan options",
		flagSet.StringVarP(&flagDir, "scan-dir", "d", "", "Directory to scan for drivers (default: Windows driver folders)"),
		flagSet.IntVarP(&opt.ScanSizeLimit, "scan-size", "l", 10, "Size limit for files to scan in MB"),
		flagSet.IntVarP(&opt.ScanWorkers, "workers", "w", 20, "Number of checksum \"threads\" to spawn"),
		flagSet.BoolVarP(&opt.ScanShowErrors, "surpress-errors", "s", false, "Do not show file read errors when calculating checksums"),
	)

	var flagGrepable bool
	var flagJson bool
	flagSet.CreateGroup("Output", "Output options",
		flagSet.BoolVarP(&flagGrepable, "grepable", "g", false, "Will only output found files for easy parsing"),
		flagSet.BoolVarP(&flagJson, "json", "j", false, "Format output as JSON"),
	)

	if err := flagSet.Parse(); err != nil {
		return nil, err
	}

	logger.Verbose = true

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
	} else if flagGrepable {
		opt.OutputMode = "grep"
		logger.Verbose = false
	} else if flagJson {
		opt.OutputMode = "json"
		logger.Verbose = false
	}

	if flagDir == "" {
		// User did not specify a path with '-d', use the default Windows directories
		opt.ScanDirectories = windowsDriverDirs
	} else {
		// User specified a custom folder to scan
		opt.ScanDirectories = append(opt.ScanDirectories, flagDir)
	}

	printBanner()

	return opt, nil
}

func printBanner() {
	logger.PlainStderr(`
  ╔─────────────────────────────────────╗
  │          LOLDrivers-client          │
  │      https://www.loldrivers.io      │
  │                                     │
  │    https://github.com/rtfmkiesel    │ 
  ╚─────────────────────────────────────╝

`)
}
