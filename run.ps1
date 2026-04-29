param(
    [switch]$run,
    [switch]$temp
)

$ErrorActionPreference = "Stop"
try { [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12 } catch {}

$repo = "rtfmkiesel/loldrivers-client"
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "arm64" }
$fileName = "LOLDrivers-client_Windows_$arch.zip"

$downloadDir = Join-Path $env:TEMP "loldrivers_download_$(Get-Random)"
New-Item -ItemType Directory -Path $downloadDir -Force | Out-Null

$extractDir = if ($temp) {
    $tempExtract = Join-Path $env:TEMP "loldrivers_$(Get-Random)"
    New-Item -ItemType Directory -Path $tempExtract -Force | Out-Null
    $tempExtract
} else {
    $PWD
}

$cleanupExtractDir = $false

try {
    Write-Host "[*] Fetching latest release..." 
    $release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
    if (-not $release) { throw "Failed to fetch release info" }
    Write-Host "[+] Found release $(($release.tag_name))" -ForegroundColor Green

    Write-Host "[*] Downloading checksum..."
    $checksumAsset = $release.assets | Where-Object { $_.name -like "loldrivers-client_*checksums.txt" }
    if (-not $checksumAsset) { throw "Checksum file not found in release" }
    
    $checksumResponse = Invoke-WebRequest $checksumAsset.browser_download_url
    if ($checksumResponse.StatusCode -ne 200) { throw "Failed to download checksum file" }
    
    $checksumContent = [System.Text.Encoding]::UTF8.GetString($checksumResponse.Content)
    if (-not $checksumContent) { throw "Checksum file is empty" }
    Write-Host "[+] Downloaded checksum" -ForegroundColor Green

    $checksumLine = $checksumContent -split "`n" | Where-Object { $_ -match $fileName }
    if (-not $checksumLine) { throw "Checksum for $fileName not found" }
    
    $expectedHash = ($checksumLine -split '\s+')[0].ToLower()
    if (-not $expectedHash) { throw "Failed to parse checksum" }

    Write-Host "[*] Downloading client..."
    $zipAsset = $release.assets | Where-Object { $_.name -eq $fileName }
    if (-not $zipAsset) { throw "$fileName not found in release" }
    
    $zipPath = Join-Path $downloadDir $fileName
    Invoke-WebRequest $zipAsset.browser_download_url -OutFile $zipPath
    
    if (-not (Test-Path $zipPath)) { throw "Failed to download $fileName" }
    Write-Host "[+] Downloaded $(($fileName))" -ForegroundColor Green
    
    Write-Host "[*] Verifying checksum..."
    $actualHash = (Get-FileHash $zipPath -Algorithm SHA256).Hash.ToLower()
    if ($actualHash -ne $expectedHash) {
        throw "Checksum mismatch: expected $expectedHash, got $actualHash"
    }
    Write-Host "[+] Checksum verified" -ForegroundColor Green

    Write-Host "[*] Extracting to $extractDir..."
    Expand-Archive $zipPath -DestinationPath $extractDir -Force
    $cleanupExtractDir = $true
    Write-Host "[+] Extracted to $extractDir\LOLDrivers-client.exe" -ForegroundColor Green

    if ($run -or $temp) {
        Write-Host "[*] Running with default arguments..."
        & "$extractDir\LOLDrivers-client.exe"
    }
}
catch {
    Write-Host "[!] ERR: $_" -ForegroundColor Red
    if ($temp -and $extractDir -and (Test-Path $extractDir)) {
        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    exit 1
}
finally {
    if ($downloadDir -and (Test-Path $downloadDir)) {
        Remove-Item $downloadDir -Recurse -Force -ErrorAction SilentlyContinue
    }

    if ($temp -and $cleanupExtractDir -and $extractDir -and (Test-Path $extractDir)) {
        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}