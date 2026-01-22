$ErrorActionPreference = "Stop"

Write-Host "Building LiteSpeed for Linux x86_64..."

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

go build -trimpath -ldflags "-s -w" -o LiteSpeed

Write-Host "Build done: ./LiteSpeed"
