# Post-Quantum Cognitive Decoys Startup Script
Write-Host "Starting PQCD Services..."

# Start the Go backend API server
Start-Process -FilePath "powershell.exe" -ArgumentList "-Command", "cd $PWD/backend && go run api_server.go" -WindowStyle Normal
Write-Host "Backend API server started at http://localhost:8082"

# Start the AI service
Start-Process -FilePath "powershell.exe" -ArgumentList "-Command", "cd $PWD/ai && python app.py" -WindowStyle Normal
Write-Host "AI service started at http://localhost:5000"

# Start the React frontend
Start-Process -FilePath "powershell.exe" -ArgumentList "-Command", "cd $PWD/frontend && npm start" -WindowStyle Normal
Write-Host "Frontend started at http://localhost:3000"

Write-Host "`nAll services started!"
Write-Host "Access the application at http://localhost:3000"
Write-Host "Backend API: http://localhost:8082"
Write-Host "AI Service API: http://localhost:5000"

# Keep this window open
Write-Host "`nPress Ctrl+C to stop all services."
while ($true) { Start-Sleep -Seconds 1 } 