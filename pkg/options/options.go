package options

import (
	"flag"
	"fmt"
	"time"

	"github.com/rtfmkiesel/loldrivers-client/pkg/logger"
)

type Options struct {
	Mode             string
	LocalDriversPath string
	ScanDirectories  []string
	ScanSizeLimit    int64
	OutputMode       string
	Workers          int
	StartTime        time.Time
}

// Parse the command line options into an Options struct
func Parse() (opt *Options, err error) {
	opt = &Options{}

	opt.StartTime = time.Now()

	var flagDir string
	var flagSilent bool
	var flagJson bool
	flag.StringVar(&opt.Mode, "m", "online", "")
	flag.StringVar(&opt.Mode, "mode", "online", "")
	flag.StringVar(&opt.LocalDriversPath, "f", "", "")
	flag.StringVar(&opt.LocalDriversPath, "driver-file", "", "")
	flag.StringVar(&flagDir, "d", "", "")
	flag.StringVar(&flagDir, "scan-dir", "", "")
	flag.Int64Var(&opt.ScanSizeLimit, "l", 10, "")
	flag.Int64Var(&opt.ScanSizeLimit, "scan-limit", 10, "")
	flag.BoolVar(&flagSilent, "s", false, "")
	flag.BoolVar(&flagSilent, "silent", false, "")
	flag.BoolVar(&flagJson, "j", false, "")
	flag.BoolVar(&flagJson, "json", false, "")
	flag.IntVar(&opt.Workers, "w", 20, "")
	flag.IntVar(&opt.Workers, "workers", 20, "")
	flag.Usage = func() { usage() }
	flag.Parse()

	switch opt.Mode {
	case "online", "internal":
		// we good
	case "local":
		if opt.Mode == "local" && opt.LocalDriversPath == "" {
			return nil, fmt.Errorf("mode 'local' requires '-f'")
		}
	default:
		return nil, fmt.Errorf("invalid mode '%s'", opt.Mode)
	}

	// Only one output style
	if flagSilent && flagJson {
		return nil, fmt.Errorf("only use '-s' or '-j', not both")
	} else if flagSilent {
		opt.OutputMode = "silent"
		logger.BeSilent = true
	} else if flagJson {
		opt.OutputMode = "json"
		logger.BeSilent = true
	}

	logger.Banner()

	// Directories
	if flagDir == "" {
		// User did not specify a path with '-d', use the default Windows opt.Directories
		opt.ScanDirectories = append(opt.ScanDirectories, "C:\\Windows\\System32\\drivers")
		opt.ScanDirectories = append(opt.ScanDirectories, "C:\\Windows\\System32\\DriverStore\\FileRepository")
		opt.ScanDirectories = append(opt.ScanDirectories, "C:\\WINDOWS\\inf")
	} else {
		// User specified a custom folder to scan
		opt.ScanDirectories = append(opt.ScanDirectories, flagDir)
	}

	return opt, nil
}

func usage() {
	logger.Banner()
	fmt.Println(`Usage:
  LOLDrivers-client.exe [OPTIONS]
 
Options:
  -m, --mode            Operating Mode {online, local, internal}
                            online = Download the newest driver set (default)
                            local = Use a local drivers.json file (requires '-f')
                            internal = Use the built-in driver set (can be outdated)

  -f, --driver-file     File path to 'drivers.json', when running in local mode

  -d, --scan-dir        Directory to scan for drivers (default: Windows driver folders)
                        Files which cannot be opened or read will be silently ignored

  -l, --scan-limit      Size limit for files to scan in MB (default: 10)
                        Be aware, higher values greatly increase runtime & CPU usage

  -s, --silent          Will only output found files for easy parsing (default: false)
  -j, --json            Format output as JSON (default: false)

  -w, --workers         Number of "threads" to spawn (default: 20)
  -h, --help            Shows this text
  `)
}
