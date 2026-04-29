<div align="center">
  <h1>LOLDrivers-client</h1>

  The OG, _blazingly fast_, **zero dependencies** client for [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) by [MagicSword](https://www.magicsword.io/).   
  Scan your Windows computer for known vulnerable or malicious drivers.

  ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rtfmkiesel/loldrivers-client/build_release.yaml) ![License](https://img.shields.io/github/license/rtfmkiesel/loldrivers-client) ![GitHub Repo stars](https://img.shields.io/github/stars/rtfmkiesel/loldrivers-client)

  ![](demo.gif)
</div>

## Installation

### Command Line

Copy and paste the command below into a PowerShell terminal. This does not require an elevated (Administrator) shell.

```ps1
# Download
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/rtfmkiesel/loldrivers-client/refs/heads/main/run.ps1)))

# Download and run
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/rtfmkiesel/loldrivers-client/refs/heads/main/run.ps1))) -run

# Download to TEMP, run, and delete
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/rtfmkiesel/loldrivers-client/refs/heads/main/run.ps1))) -temp
```

### Download

Download the prebuilt binaries [from GitHub](https://github.com/rtfmkiesel/loldrivers-client/releases).

### Build From Source

```sh
# requires Golang >=1.26

git clone https://github.com/rtfmkiesel/loldrivers-client
cd loldrivers-client
go generate ./internal/loldrivers/
go test ./internal/loldrivers/
go build -o LOLDrivers-client.exe -ldflags="-s -w" .
```
## Usage

```
Usage of LOLDrivers-client.exe:
  -debug
        print debug output (will mess up 'grep' and 'json' output)
  -file string
        /path/to/drivers.json (only for '-mode file')
  -maxsize int
        size limit for files to scan in MB (default 10)
  -mode string
        operating mode {online, file, embedded} (default "online")
  -output string
        output mode {standard,grep,json} (default "standard")
  -target string
        target directory (default=OS)
  -workers int
        number of parallel scan workers (default 20)
```

## Legal

This project is not affiliated with the [LOLDrivers](https://github.com/magicsword-io/LOLDrivers) project.
