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
	flag.BoolVar(&logger.NoColor, "nocolor", false, "do not print colored output")
	target := flag.String("target", "", "target directory (default OS)")
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
	case "standard":
		c.outputmode = *output
	case "grep", "json":
		logger.Silent = true
		c.outputmode = *output
	default:
		return nil, fmt.Errorf("config: invalid '-output'")
	}

	if c.outputmode == "standard" {
		logger.Stderr(`
┌─────────────────────────────────────────────────────────────────────┐
│ ____   ___________   _____  ______________ ________________________ │
│ 7  7   7     77  7   7    \ 7  _  77  77  V  77     77  _  77     7 │
│ |  |   |  7  ||  |   |  7  ||    _||  ||  |  ||  ___!|    _||  ___! │
│ |  !___|  |  ||  !___|  |  ||  _ \ |  ||  !  ||  __|_|  _ \ !__   7 │
│ |     7|  !  ||     7|  !  ||  7  ||  ||     ||     7|  7  |7     | │
│ !_____!!_____!!_____!!_____!!__!__!!__!!_____!!_____!!__!__!!_____! │
│                                                                     │
│           https://github.com/rtfmkiesel/loldrivers-client           │
└─────────────────────────────────────────────────────────────────────┘

`)
	}

	return c, nil
}
