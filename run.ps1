# Post-Quantum Cognitive Decoys System Startup Script

Write-Host "Starting Post-Quantum Cognitive Decoys System" -ForegroundColor Cyan

# Initialize the database
Write-Host "Initializing database..." -ForegroundColor Green
Push-Location -Path "database"
python init_db.py
Pop-Location

# Start the AI service in a new window
Write-Host "Starting AI service..." -ForegroundColor Green
$aiProcess = Start-Process -FilePath "python" -ArgumentList "ai/app.py" -PassThru -WindowStyle Normal

# Start the Go backend in a new window
Write-Host "Starting Go backend..." -ForegroundColor Green
$backendProcess = Start-Process -FilePath "go" -ArgumentList "run", "backend/simple.go" -PassThru -WindowStyle Normal

Write-Host "All services started." -ForegroundColor Cyan
Write-Host "Backend API: http://localhost:8083" -ForegroundColor Yellow
Write-Host "AI Service: http://localhost:5000" -ForegroundColor Yellow
Write-Host "Press Ctrl+C in each window to stop the services." -ForegroundColor Yellow

# Register cleanup function
$null = Register-EngineEvent -SourceIdentifier ([System.Management.Automation.PsEngineEvent]::Exiting) -Action {
    Write-Host "Shutting down services..." -ForegroundColor Red
    if ($aiProcess -ne $null -and -not $aiProcess.HasExited) {
        Stop-Process -Id $aiProcess.Id -Force
    }
    if ($backendProcess -ne $null -and -not $backendProcess.HasExited) {
        Stop-Process -Id $backendProcess.Id -Force
    }
} 