$ErrorActionPreference = "Stop"

Write-Host "Building LiteSpeed for Linux x86_64..."

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

go build -trimpath -ldflags "-s -w" -o lite-linux-with-ip

Write-Host "Build done: ./lite-linux-with-ip"
