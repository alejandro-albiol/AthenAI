# AthenAI Super Admin Setup Script
# Use this script to create platform super administrators
# WARNING: Only run this in secure environments with proper access controls

param(
    [switch]$Force,
    [string]$Environment = "development"
)

Write-Host "=== AthenAI Super Admin Setup ===" -ForegroundColor Yellow
Write-Host "Environment: $Environment" -ForegroundColor Cyan
Write-Host ""

# Security check
if ($Environment -eq "production" -and -not $Force) {
    Write-Host "ERROR: Production environment detected!" -ForegroundColor Red
    Write-Host "Creating super admins in production requires explicit --Force flag" -ForegroundColor Red
    Write-Host "Ensure you have proper security approvals before proceeding." -ForegroundColor Red
    exit 1
}

# Check if Go is available
$goVersion = go version 2>$null
if (-not $goVersion) {
    Write-Host "ERROR: Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

Write-Host "Go version: $goVersion" -ForegroundColor Green

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

# Final confirmation
if ($Environment -eq "production") {
    Write-Host "⚠️  WARNING: You are about to create a super admin in PRODUCTION!" -ForegroundColor Red
    Write-Host "This will give someone complete control over the platform." -ForegroundColor Red
    Write-Host ""
    $confirm = Read-Host "Type 'CREATE-PRODUCTION-SUPERADMIN' to continue"
    if ($confirm -ne "CREATE-PRODUCTION-SUPERADMIN") {
        Write-Host "Operation cancelled." -ForegroundColor Yellow
        exit 0
    }
}

# Run the super admin setup
Write-Host "Starting super admin creation process..." -ForegroundColor Yellow
Write-Host ""

go run cmd/setup-superadmin/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "=== SECURITY REMINDERS ===" -ForegroundColor Red
    Write-Host "1. Store credentials in a secure password manager" -ForegroundColor White
    Write-Host "2. Enable 2FA for the admin account (when implemented)" -ForegroundColor White
    Write-Host "3. Monitor admin access logs regularly" -ForegroundColor White
    Write-Host "4. Remove this setup tool from production servers" -ForegroundColor White
    Write-Host "5. Document who has super admin access" -ForegroundColor White
    Write-Host ""
    Write-Host "Super admin setup completed successfully!" -ForegroundColor Green
}
else {
    Write-Host "Super admin setup failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}
