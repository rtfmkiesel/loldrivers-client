# LOLDrivers-client
![GitHub Repo stars](https://img.shields.io/github/stars/rtfmkiesel/loldrivers-client) ![GitHub](https://img.shields.io/github/license/rtfmkiesel/loldrivers-client)

The first *blazingly fast* client for [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) (Living Off The Land Drivers) by [MagicSword](https://www.magicsword.io/). Scan your computer for known vulnerable and known malicious Windows drivers.


![](demo.gif)


## Usage
```
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
```

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

# Legal
This project is not affiliated with the [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) repository.