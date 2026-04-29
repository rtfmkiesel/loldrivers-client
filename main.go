//go:build windows

package main

import (
	"github.com/rtfmkiesel/loldrivers-client/cmd"
	"github.com/rtfmkiesel/loldrivers-client/internal/logger"
)

func main() {
	if err := cmd.Run(); err != nil {
		logger.Fatal(err)
	}
}
