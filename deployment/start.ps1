# PowerShell script to start the PQCD application using Docker Compose

Write-Host "Checking if Docker is running..." -ForegroundColor Cyan

try {
    $dockerStatus = docker info 2>$null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Docker is not running. Please start Docker Desktop first." -ForegroundColor Red
        exit 1
    }
    Write-Host "Docker is running." -ForegroundColor Green
} catch {
    Write-Host "Failed to check Docker status: $_" -ForegroundColor Red
    Write-Host "Please ensure Docker is installed and running." -ForegroundColor Red
    exit 1
}

Write-Host "Starting the PQCD application..." -ForegroundColor Cyan
Write-Host "This will build and start all services: backend, AI service, and frontend." -ForegroundColor Cyan

# Build and start services
try {
    docker-compose -f docker-compose.yml up --build
} catch {
    Write-Host "Failed to start Docker Compose: $_" -ForegroundColor Red
    exit 1
} 