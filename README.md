# LOLDrivers-client
![GitHub Repo stars](https://img.shields.io/github/stars/rtfmkiesel/loldrivers-client) ![GitHub](https://img.shields.io/github/license/rtfmkiesel/loldrivers-client)

The first *blazingly fast* client for [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) (Living Off The Land Drivers) by [MagicSword](https://www.magicsword.io/). Scan your computer for known vulnerable and known malicious Windows drivers.

![](demo.gif)

## Usage
```
Usage:
  loldrivers-client.exe [flags]

Flags:
OPERATING MODE:
   -m, -mode string         Operating Mode {online, local, internal} (default "online")
   -f, -driver-file string  File path to 'drivers.json', when mode == local

SCAN OPTIONS:
   -d, -scan-dir string  Directory to scan for drivers (default: Windows driver folders)
   -l, -scan-size int    Size limit for files to scan in MB (default 10)
   -w, -workers int      Number of checksum "threads" to spawn (default 20)
   -s, -surpress-errors  Do not show file read errors when calculating checksums

OUTPUT OPTIONS:
   -g, -grepable  Will only output found files for easy parsing
   -j, -json      Format output as JSON
```

## Installation
### Binaries
Download the prebuilt binaries [here](https://github.com/rtfmkiesel/loldrivers-client/releases).

## Build from source
```
git clone https://github.com/rtfmkiesel/loldrivers-client
cd loldrivers-client
go mod tidy
go build -o LOLDrivers-client.exe -ldflags="-s -w" cmd/loldrivers-client/loldrivers-client.go
```

# Contributing 
Improvements in the form of PRs are always welcome, especially as this was made during my first year of using Golang. 

# Legal
This project is not affiliated with the [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) repository.
