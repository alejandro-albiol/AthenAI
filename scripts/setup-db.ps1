# Database Setup Script
Write-Host "=== Athenai Database Setup ===" -ForegroundColor Cyan
Write-Host ""

# Check if .env file exists
if (!(Test-Path ".env")) {
    Write-Host "Error: .env file not found!" -ForegroundColor Red
    Write-Host "Please create a .env file with your database configuration." -ForegroundColor Yellow
    Write-Host "You can copy example.env and update the values." -ForegroundColor Yellow
    exit 1
}

Write-Host "Loading environment variables..." -ForegroundColor Yellow
dotenv -f .env

Write-Host "Connecting to database..." -ForegroundColor Yellow

# Build and run the database setup
Write-Host "Building database setup..." -ForegroundColor Yellow
go build -o tmp/db-setup cmd/setup-db/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Running database setup..." -ForegroundColor Yellow
./tmp/db-setup

if ($LASTEXITCODE -ne 0) {
    Write-Host "Database setup failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Database setup completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Tables created:" -ForegroundColor Cyan
Write-Host "  - public.gym" -ForegroundColor White
Write-Host "  - public.admin" -ForegroundColor White
Write-Host "  - public.exercise" -ForegroundColor White
Write-Host "  - public.muscular_group" -ForegroundColor White
Write-Host "  - public.equipment" -ForegroundColor White
Write-Host ""
Write-Host "All indexes have been created for optimal performance." -ForegroundColor Green 