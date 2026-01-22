$ErrorActionPreference = "Stop"

# ----------------------------
# Config
# ----------------------------
$AppName = "LiteSpeed"
$OutDir  = "dist"
$MainPkg = "."   # 若你的 main 在子目录，改成 "./cmd/LiteSpeed"

Write-Host "Building $AppName"
Write-Host "Main package: $MainPkg"

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

# ----------------------------
# Targets
# ----------------------------
# Linux:
# - amd64 = x86_64
# - 386   = x86
# - arm64 = aarch64
# - armv6/armv7 via GOARCH=arm + GOARM=6/7
#
# Windows:
# - amd64 (x86_64)
# - 386   (x86)
# - arm64 (Windows on ARM)
$Targets = @(
    # Linux
    @{ GOOS="linux";   GOARCH="amd64"; GOARM="";  Suffix="amd64"; Ext=""    },
    @{ GOOS="linux";   GOARCH="386";   GOARM="";  Suffix="386";   Ext=""    },
    @{ GOOS="linux";   GOARCH="arm64"; GOARM="";  Suffix="arm64"; Ext=""    },
    @{ GOOS="linux";   GOARCH="arm";   GOARM="6"; Suffix="armv6"; Ext=""    },
    @{ GOOS="linux";   GOARCH="arm";   GOARM="7"; Suffix="armv7"; Ext=""    },

    # Windows
    @{ GOOS="windows"; GOARCH="amd64"; GOARM="";  Suffix="amd64"; Ext=".exe"},
    @{ GOOS="windows"; GOARCH="386";   GOARM="";  Suffix="386";   Ext=".exe"},
    @{ GOOS="windows"; GOARCH="arm64"; GOARM="";  Suffix="arm64"; Ext=".exe"}
)

# Common build flags
$TrimPath = "-trimpath"
$LdFlags  = "-s -w"

Write-Host "Output directory: $OutDir"
Write-Host ""

foreach ($t in $Targets) {
    $goos   = $t.GOOS
    $goarch = $t.GOARCH
    $goarm  = $t.GOARM
    $suffix = $t.Suffix
    $ext    = $t.Ext

    # Set env for Go
    $env:GOOS = $goos
    $env:GOARCH = $goarch
    $env:CGO_ENABLED = "0"

    if ($goarch -eq "arm") {
        $env:GOARM = $goarm
    } else {
        Remove-Item Env:\GOARM -ErrorAction SilentlyContinue
    }

    $outName = "{0}_{1}_{2}{3}" -f $AppName, $goos, $suffix, $ext
    $outPath = Join-Path $OutDir $outName

    $armNote = ""
    if ($goarch -eq "arm" -and $goarm) { $armNote = " (GOARM=$goarm)" }

    Write-Host ("Building {0}/{1}{2} -> {3}" -f $goos, $goarch, $armNote, $outPath)

    go build $TrimPath -ldflags $LdFlags -o $outPath $MainPkg

    Write-Host "Done."
    Write-Host ""
}

Write-Host "All builds completed. Artifacts in ./$OutDir"
