#!/usr/bin/env pwsh
# PowerShell script to start all PQCD components for local development

# Create function to start a process in a new window
function Start-ProcessInNewWindow {
    param (
        [string]$Title,
        [string]$Command,
        [string]$WorkingDirectory
    )
    
    Write-Host "Starting $Title..."
    Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$WorkingDirectory'; Write-Host 'Starting $Title...'; $Command" -WindowStyle Normal
}

# Initialize the database if it doesn't exist
if (-not (Test-Path "./database/pqcd.db")) {
    Write-Host "Initializing database..."
    Push-Location "./database"
    python init_db.py "./pqcd.db"
    Pop-Location
}

# Start the backend
Start-ProcessInNewWindow -Title "PQCD Backend" -Command "cd ./backend; go run main.go" -WorkingDirectory (Get-Location)

# Start the AI service
Start-ProcessInNewWindow -Title "PQCD AI Service" -Command "cd ./ai; python app.py" -WorkingDirectory (Get-Location)

# Start the frontend
Start-ProcessInNewWindow -Title "PQCD Frontend" -Command "cd ./frontend; npm start" -WorkingDirectory (Get-Location)

Write-Host "All services started. Access the application at http://localhost:3000"
Write-Host "Backend API: http://localhost:8082"
Write-Host "AI Service API: http://localhost:5000" 