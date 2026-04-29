package cmd

import (
	"flag"
	"fmt"

	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
)

type config struct {
	mode        string
	localFile   string
	targets     []string
	maxSize     int64
	workerCount int
	outputmode  string
}

func parseConfig() (*config, error) {
	c := &config{}

	flag.StringVar(&c.mode, "mode", "online", "operating mode {online, file, embedded}")
	flag.StringVar(&c.localFile, "file", "", "/path/to/drivers.json (only for '-mode file')")
	flag.IntVar(&c.workerCount, "workers", 20, "number of parallel scan workers")
	target := flag.String("target", "", "target directory (default=OS)")
	limitMb := flag.Int64("maxsize", 10, "size limit for files to scan in MB")
	output := flag.String("output", "standard", "output mode {standard,grep,json}")
	debug := flag.Bool("debug", false, "print debug output (will mess up 'grep' and 'json' output)")
	flag.Parse()

	logger.PrintDebug = *debug

	switch c.mode {
	case "online", "embedded":
		// we good
	case "file":
		if c.localFile == "" {
			return nil, fmt.Errorf("config: '-mode file' requires '-file /path/to/drivers.json'")
		}
	default:
		return nil, fmt.Errorf("config: invalid '-mode'")
	}

	if *target != "" {
		c.targets = []string{*target} // User-defined
	} else {
		// Default OS
		c.targets = []string{
			"C:\\Windows\\System32\\drivers",
			"C:\\Windows\\System32\\DriverStore\\FileRepository",
			"C:\\WINDOWS\\inf",
		}
	}

	c.maxSize = *limitMb * 1024 * 1024 // MB to bytes

	switch *output {
	case "standard", "grep", "json":
		c.outputmode = *output
	default:
		return nil, fmt.Errorf("config: invalid '-output'")
	}

	if c.outputmode == "standard" {
		logger.Stderr(`
╔─────────────────────────────────────╗
│    LOLDrivers-client                │
│    https://www.loldrivers.io        │
|                                     |
|    made by                          |
|    https://github.com/rtfmkiesel    |
╚─────────────────────────────────────╝

`)
	}

	return c, nil
}
