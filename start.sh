#!/bin/bash
# Bash script to start all PQCD components for local development

# Function to start a process in a new terminal
start_in_terminal() {
    local title=$1
    local command=$2
    local dir=$3
    
    echo "Starting $title..."
    
    if command -v gnome-terminal &> /dev/null; then
        gnome-terminal --title="$title" -- bash -c "cd '$dir' && $command; exec bash"
    elif command -v xterm &> /dev/null; then
        xterm -title "$title" -e "cd '$dir' && $command; exec bash" &
    elif command -v konsole &> /dev/null; then
        konsole --new-tab -p tabtitle="$title" -e bash -c "cd '$dir' && $command; exec bash" &
    elif command -v terminal &> /dev/null; then
        terminal -t "$title" -e "cd '$dir' && $command; exec bash" &
    else
        echo "No supported terminal found. Running $title in background..."
        (cd "$dir" && $command) &
    fi
}

# Initialize the database if it doesn't exist
if [ ! -f "./database/pqcd.db" ]; then
    echo "Initializing database..."
    (cd ./database && python init_db.py "./pqcd.db")
fi

# Start the backend
start_in_terminal "PQCD Backend" "go run main.go" "$(pwd)/backend"

# Start the AI service
start_in_terminal "PQCD AI Service" "python app.py" "$(pwd)/ai"

# Start the frontend
start_in_terminal "PQCD Frontend" "npm start" "$(pwd)/frontend"

echo "All services started. Access the application at http://localhost:3000"
echo "Backend API: http://localhost:8082"
echo "AI Service API: http://localhost:5000" 