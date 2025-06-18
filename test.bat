@echo off
echo Starting Post-Quantum Cognitive Decoys System for testing...

REM Start the backend service
start "PQCD Backend" cmd /c "cd backend && go run simple.go"
echo Backend service starting on port 8083...

REM Start the AI service
start "PQCD AI Service" cmd /c "cd ai && python app.py"
echo AI service starting on port 5000...

REM Wait for services to start
echo Waiting for services to start...
timeout /t 5 /nobreak

REM Run the test script
echo Running tests...
python test_system.py

REM Keep the window open to see results
pause 