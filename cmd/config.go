package cmd

import (
	"flag"
	"fmt"

	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
)

type config struct {
	targets     []string
	maxSize     int64
	workerCount int
	outputmode  string
}

func parseConfig() (*config, error) {
	c := &config{}
	flag.IntVar(&c.workerCount, "workers", 20, "number of parallel scan workers")
	flag.BoolVar(&logger.PrintDebug, "debug", false, "print debug output (will mess up 'grep' and 'json' output)")
	flag.BoolVar(&logger.NoColor, "nocolor", false, "do not print colored output")
	target := flag.String("target", "", "target directory (default=OS)")
	limitMb := flag.Int64("maxsize", 10, "size limit for files to scan in MB")
	output := flag.String("output", "standard", "output mode {standard,grep,json}")
	flag.Parse()

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
