#!/bin/bash

# Set up error handling
set -e

echo "Starting Post-Quantum Cognitive Decoys System"

# Initialize the database
echo "Initializing database..."
cd database
python init_db.py
cd ..

# Start the AI service in the background
echo "Starting AI service..."
cd ai
python app.py &
AI_PID=$!
cd ..

# Start the Go backend
echo "Starting Go backend..."
cd backend
go run simple.go &
BACKEND_PID=$!
cd ..

# Trap to handle Ctrl+C and cleanup
trap 'echo "Shutting down services..."; kill $AI_PID $BACKEND_PID; exit' INT TERM

echo "All services started. Press Ctrl+C to stop."
echo "Backend API: http://localhost:8083"
echo "AI Service: http://localhost:5000"

# Wait for processes to finish
wait $AI_PID $BACKEND_PID 