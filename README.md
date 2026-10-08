<div align="center">
  <h1>LOLDrivers-client</h1>

  The OG, _blazingly fast_, **zero dependencies** client for [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) by [MagicSword](https://www.magicsword.io/).   
  Scan your Windows computer for known vulnerable or malicious drivers.

  ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rtfmkiesel/loldrivers-client/build_release.yaml) ![License](https://img.shields.io/github/license/rtfmkiesel/loldrivers-client) ![GitHub Repo stars](https://img.shields.io/github/stars/rtfmkiesel/loldrivers-client)

  ![](demo.gif)
</div>

## Installation

### Command Line

Copy and paste the command below into a PowerShell terminal. **This does not require an elevated (Administrator) shell.**

```ps1
# Downloads the client to the current directory
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/rtfmkiesel/loldrivers-client/refs/heads/main/run.ps1)))

# Downloads, runs and deletes the client (temp, default options)
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/rtfmkiesel/loldrivers-client/refs/heads/main/run.ps1))) -temp
```

### Download

Download the prebuilt binaries [from GitHub](https://github.com/rtfmkiesel/loldrivers-client/releases).  

The `*_embedded.zip` version has an embedded [drivers.json](https://www.loldrivers.io/api/drivers.json) file for offline usability. The age of the JSON file is determined by the last build/release date. Please note, this file [might get flagged as malware](https://github.com/rtfmkiesel/loldrivers-client/issues/4).

### Build From Source

```sh
# Requires Golang >=1.27 and a Linux env
# When using Windows, change 'curl' to 'curl.exe' in 'internal\loldrivers\drivers_embedded.go'

git clone https://github.com/rtfmkiesel/loldrivers-client
cd loldrivers-client

go generate ./internal/loldrivers/

# normal
go build -o LOLDrivers-client.exe -ldflags="-s -w" .
# embedded
go build -o LOLDrivers-client_embedded.exe -ldflags="-s -w" -tags embedded .
```

## Usage

```
Usage of LOLDrivers-client.exe:
  -debug
        print debug output (will mess up 'grep' and 'json' output)
  -maxsize int
        size limit for files to scan in MB (default 10)
  -nocolor
        do not print colored output
  -output string
        output mode {standard,grep,json} (default "standard")
  -target string
        target directory (default=OS)
  -workers int
        number of parallel scan workers (default 20)
```

## Legal

This project is not affiliated with the [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) project.
