<#
.SYNOPSIS
    ContextKeeper Installation Script for Windows

.DESCRIPTION
    Installs ContextKeeper CLI tool on Windows systems

.PARAMETER Version
    Version to install (default: latest)

.PARAMETER InstallDir
    Installation directory (default: C:\Program Files\ContextKeeper)
#>

param(
    [string]$Version = "latest",
    [string]$InstallDir = "C:\Program Files\ContextKeeper"
)

$ErrorActionPreference = "Stop"

function Write-Info {
    param([string]$Message)
    Write-Host "INFO: $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "WARNING: $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "ERROR: $Message" -ForegroundColor Red
}

function Get-LatestVersion {
    if ($Version -eq "latest") {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/ondrahracek/contextkeeper/releases/latest"
        return $response.tag_name
    }
    return $Version
}

function Install-Binary {
    param(
        [string]$Version,
        [string]$InstallDir
    )
    
    $os = "windows"
    $arch = "amd64"
    $binaryName = "contextkeeper.exe"
    $url = "https://github.com/ondrahracek/contextkeeper/releases/download/$Version/contextkeeper-$os-$arch.exe"
    $tmpDir = [System.IO.Path]::GetTempFileName()
    
    Write-Info "Downloading ContextKeeper $Version for $os-$arch..."
    
    try {
        Invoke-WebRequest -Uri $url -OutFile $tmpDir -UseBasicParsing
        
        if (-not (Test-Path $tmpDir)) {
            throw "Download failed"
        }
        
        New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
        Copy-Item $tmpDir "$InstallDir\$binaryName"
        
        # Add to PATH (for current session)
        $env:PATH = "$InstallDir;$env:PATH"
        
        Write-Info "Installed to $InstallDir\$binaryName"
        
        # Cleanup
        Remove-Item $tmpDir -ErrorAction SilentlyContinue
        
    } catch {
        Write-Error "Download failed: $_"
        throw
    }
}

function Verify-Install {
    try {
        Write-Info "Verifying installation..."
        & "$InstallDir\contextkeeper.exe" --version
        Write-Info "Installation successful!"
    } catch {
        Write-Warn "Verification failed. Binary may not be in PATH."
    }
}

function Add-ToSystemPath {
    param([string]$Path)
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    if ($currentPath -notlike "*$Path*") {
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$Path", "Machine")
        Write-Info "Added $Path to system PATH"
        Write-Warn "You may need to restart your terminal for changes to take effect"
    }
}

# Main
Write-Host "==============================================" -ForegroundColor Cyan
Write-Host "  ContextKeeper Installation Script (Windows)" -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Cyan
Write-Host ""

$version = Get-LatestVersion
Write-Host "Version: $version" -ForegroundColor White
Write-Host "Install Directory: $InstallDir" -ForegroundColor White
Write-Host ""

Install-Binary -Version $version -InstallDir $InstallDir
Verify-Install

Write-Host ""
Write-Info "Usage: contextkeeper add 'Your context note'"
Write-Info "       contextkeeper list"
Write-Host ""
