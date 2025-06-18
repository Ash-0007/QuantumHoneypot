#!/usr/bin/env python3
import sqlite3
import os
import sys
import logging
import time
from datetime import datetime

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("db-init")

def init_db(db_path):
    """Initialize the database with schema and seed data"""
    # Check if the database directory exists
    db_dir = os.path.dirname(db_path)
    if db_dir and not os.path.exists(db_dir):
        os.makedirs(db_dir)
        logger.info(f"Created database directory: {db_dir}")

    # Wait for the directory to be accessible in Docker environment
    retries = 5
    while retries > 0 and not os.path.exists(db_dir):
        logger.info(f"Waiting for database directory to be available: {db_dir}")
        time.sleep(1)
        retries -= 1

    # Connect to the database (creates it if it doesn't exist)
    logger.info(f"Initializing database at {db_path}")
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    try:
        # Read schema from file
        with open('schema.sql', 'r') as f:
            schema_sql = f.read()

        # Execute the schema
        cursor.executescript(schema_sql)
        conn.commit()
        logger.info("Database schema created successfully")

        # Check if tables need to be seeded
        cursor.execute("SELECT COUNT(*) FROM key_pairs")
        if cursor.fetchone()[0] == 0:
            seed_database(conn)
            logger.info("Database seeded with initial data")
        else:
            logger.info("Database already contains data, skipping seeding")

    except sqlite3.Error as e:
        logger.error(f"Database initialization error: {e}")
        conn.close()
        sys.exit(1)

    conn.close()
    logger.info("Database initialization completed successfully")

def seed_database(conn):
    """Seed the database with initial data"""
    cursor = conn.cursor()

    # Sample algorithms for post-quantum cryptography
    pq_algorithms = [
        "kyber512", "kyber768", "kyber1024",
        "saber", "lightsaber", "firesaber",
        "ntru-hps-2048-509", "ntru-hps-4096-821",
        "dilithium2", "dilithium3", "dilithium5",
        "falcon512", "falcon1024",
        "sphincs-haraka-128f", "sphincs-haraka-256f",
        "frodokem-640-aes", "frodokem-976-aes", "frodokem-1344-aes"
    ]

    # Insert sample decoys
    for i, algo in enumerate(pq_algorithms):
        target = algo
        decoy = f"{algo}-variant-{i}"
        complexity = (i % 10) + 1
        
        cursor.execute(
            "INSERT INTO decoys (decoy_text, target_text, complexity, effectiveness_score) VALUES (?, ?, ?, ?)",
            (decoy, target, complexity, 0.7 + (i * 0.01))
        )

    # Insert sample event logs
    events = [
        ("system", "Database initialized", "INFO"),
        ("crypto", "Generated test key pairs", "INFO"),
        ("security", "Detected invalid access attempt", "WARNING"),
        ("system", "AI model loaded successfully", "INFO"),
        ("crypto", "Key rotation performed", "INFO")
    ]

    for event_type, description, severity in events:
        cursor.execute(
            "INSERT INTO event_logs (event_type, description, severity, timestamp) VALUES (?, ?, ?, ?)",
            (event_type, description, severity, datetime.now().isoformat())
        )

    conn.commit()

if __name__ == "__main__":
    # Get database path from command line argument or use default
    db_path = sys.argv[1] if len(sys.argv) > 1 else "./pqcd.db"
    init_db(db_path) 