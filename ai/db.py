import sqlite3
import os
import logging

logger = logging.getLogger("pqcd-ai-db")

def get_db_connection(db_path=None):
    """Get a database connection with proper configuration"""
    if db_path is None:
        db_path = os.environ.get('DB_PATH', './pqcd.db')
    
    try:
        conn = sqlite3.connect(db_path)
        conn.row_factory = sqlite3.Row  # Return rows as dictionaries
        return conn
    except sqlite3.Error as e:
        logger.error(f"Database connection error: {e}")
        raise

def store_decoy(decoy_text, target_text, complexity, effectiveness_score=None):
    """Store a generated decoy in the database"""
    db_path = os.environ.get('DB_PATH', './pqcd.db')
    
    try:
        conn = get_db_connection(db_path)
        cursor = conn.cursor()
        
        cursor.execute(
            "INSERT INTO decoys (decoy_text, target_text, complexity, effectiveness_score) VALUES (?, ?, ?, ?)",
            (decoy_text, target_text, complexity, effectiveness_score)
        )
        
        conn.commit()
        conn.close()
        return True
    except sqlite3.Error as e:
        logger.error(f"Failed to store decoy: {e}")
        return False

def log_event(event_type, description, severity="INFO"):
    """Log an event to the database"""
    db_path = os.environ.get('DB_PATH', './pqcd.db')
    
    try:
        conn = get_db_connection(db_path)
        cursor = conn.cursor()
        
        cursor.execute(
            "INSERT INTO event_logs (event_type, description, severity) VALUES (?, ?, ?)",
            (event_type, description, severity)
        )
        
        conn.commit()
        conn.close()
        return True
    except sqlite3.Error as e:
        logger.error(f"Failed to log event: {e}")
        return False

def get_stored_decoys(target_text=None, limit=100):
    """Retrieve stored decoys from the database"""
    db_path = os.environ.get('DB_PATH', './pqcd.db')
    
    try:
        conn = get_db_connection(db_path)
        cursor = conn.cursor()
        
        if target_text:
            cursor.execute(
                "SELECT * FROM decoys WHERE target_text = ? ORDER BY created_at DESC LIMIT ?",
                (target_text, limit)
            )
        else:
            cursor.execute(
                "SELECT * FROM decoys ORDER BY created_at DESC LIMIT ?",
                (limit,)
            )
            
        rows = cursor.fetchall()
        decoys = [dict(row) for row in rows]
        
        conn.close()
        return decoys
    except sqlite3.Error as e:
        logger.error(f"Failed to retrieve decoys: {e}")
        return [] 