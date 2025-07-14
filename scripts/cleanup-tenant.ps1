# AthenAI Tenant Cleanup Script
# Use this script to clean up tenant schemas for deleted gyms

param(
    [string]$SchemaName = "",
    [switch]$List,
    [switch]$Force
)

Write-Host "=== AthenAI Tenant Schema Cleanup ===" -ForegroundColor Yellow
Write-Host ""

# Check if Go is available
$goVersion = go version 2>$null
if (-not $goVersion) {
    Write-Host "ERROR: Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

# Check if database is accessible
Write-Host "Checking database connection..." -ForegroundColor Yellow
go run cmd/setup-db/main.go --test-connection 2>&1 | Out-Null
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Cannot connect to database" -ForegroundColor Red
    Write-Host "Make sure your database is running and environment variables are set" -ForegroundColor Red
    exit 1
}

Write-Host "Database connection: OK" -ForegroundColor Green
Write-Host ""

if ($List) {
    Write-Host "📋 Listing all tenant schemas..." -ForegroundColor Cyan
    Write-Host ""
    go run cmd/cleanup-tenant/main.go
    exit 0
}

if ($SchemaName -ne "") {
    if (-not $Force) {
        Write-Host "⚠️  You specified schema: $SchemaName" -ForegroundColor Yellow
        Write-Host "This will permanently delete all data in this schema." -ForegroundColor Red
        Write-Host "Use -Force flag if you're sure, or run without -SchemaName to use interactive mode." -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "🗑️  Deleting schema: $SchemaName (forced)" -ForegroundColor Red
    $env:DELETE_SCHEMA = $SchemaName
    go run cmd/cleanup-tenant/main.go
}
else {
    Write-Host "🔧 Running interactive tenant cleanup..." -ForegroundColor Cyan
    Write-Host ""
    go run cmd/cleanup-tenant/main.go
}

Write-Host ""
Write-Host "=== Cleanup Complete ===" -ForegroundColor Green
