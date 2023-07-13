# LOLDrivers-client
A client for [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) (Living Off The Land Drivers). Scan your computer for known vulnerable and known malicious Windows drivers.

## Usage
```
LOLDrivers-client.exe -m [MODE] [OPTIONS]

Modes:
  online    Download the newest driver set (default)
  local     Use a local drivers.json file (requires '-f')
  internal  Use the built-in driver set (can be outdated, fallback)

Options:
  -d        Directory to scan for drivers (default: Windows Default Driver Folders)
            Files which cannot be opened or read will be silently ignored
  -l        Size limit for files to scan in MB (default: 10)
            Be aware, higher values greatly increase runtime & CPU usage
  -f        File path to 'drivers.json'
            Only needed with '-m local'
  -t        Number of threads to spawn (default: 20)
  -v        Print verbose messages (default: false)
  -h        Shows this text
```
**Warning:** This project is not affiliated with the [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) repository. JSON structure changes in the LOLDrivers API may break this client. Since this client gets compiled with a working data set, it will fall back to the internal data set, if the parsing of the online data or the local file was not successful.

## Installation
### Binaries
Download the prebuilt binaries [here](https://github.com/rtfmkiesel/loldrivers-client/releases).

## Build from source
```bash
git clone https://github.com/rtfmkiesel/loldrivers-client
cd loldrivers-client
go mod tidy
go build -o LOLDrivers-client.exe -ldflags="-s -w" cli/loldrivers-client/loldrivers-client.go
```

# Contributing 
Improvements in the form of PRs are always welcome, especially as this was made during my first year of using Golang. 