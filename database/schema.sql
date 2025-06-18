-- Schema for Post-Quantum Cognitive Decoys database

-- Key pair table
CREATE TABLE IF NOT EXISTS key_pairs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    public_key BLOB NOT NULL,
    private_key BLOB NOT NULL,
    fingerprint TEXT NOT NULL,
    algorithm TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_real BOOLEAN DEFAULT 1,
    tags TEXT
);

-- Create index on fingerprint for faster lookups
CREATE INDEX IF NOT EXISTS idx_key_pairs_fingerprint ON key_pairs(fingerprint);

-- Create index on algorithm type
CREATE INDEX IF NOT EXISTS idx_key_pairs_algorithm ON key_pairs(algorithm);

-- Decoy table 
CREATE TABLE IF NOT EXISTS decoys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    decoy_text TEXT NOT NULL,
    target_text TEXT NOT NULL,
    complexity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    effectiveness_score REAL
);

-- Create index on target text
CREATE INDEX IF NOT EXISTS idx_decoys_target ON decoys(target_text);

-- Event log table
CREATE TABLE IF NOT EXISTS event_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL,
    description TEXT,
    source_ip TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    severity TEXT CHECK (severity IN ('INFO', 'WARNING', 'ERROR', 'CRITICAL')),
    related_item_id INTEGER,
    related_item_type TEXT
);

-- Create index on event type and timestamp
CREATE INDEX IF NOT EXISTS idx_event_logs_type_time ON event_logs(event_type, timestamp);

-- User table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,  -- Store hashed passwords only
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    role TEXT CHECK (role IN ('admin', 'user', 'readonly')) DEFAULT 'user'
);

-- Insert default admin user (password: admin123)
INSERT OR IGNORE INTO users (username, password_hash, role) 
VALUES ('admin', '$2b$12$1B7tQMt1IjuOXZvGdyz9A.J1DWnbwOgqYJMBpWY2Rcx.UmNMy9.cG', 'admin'); 