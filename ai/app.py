from flask import Flask, request, jsonify
from flask_cors import CORS
import numpy as np
import pickle
import os
import time
import logging
import sqlite3
from sklearn.ensemble import RandomForestClassifier
from sklearn.feature_extraction.text import CountVectorizer
from datetime import datetime
import re
from difflib import SequenceMatcher
import random
import string

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("pqcd-ai")

app = Flask(__name__)
CORS(app)  # Add CORS support for all routes

# Model will be loaded on first request
model = None
vectorizer = None

def load_or_create_model():
    """Load existing model or create a new one if it doesn't exist"""
    global model, vectorizer
    
    # Create models directory if it doesn't exist
    os.makedirs("models", exist_ok=True)
    
    model_path = os.path.join("models", "decoy_model.pkl")
    vectorizer_path = os.path.join("models", "vectorizer.pkl")

    if os.path.exists(model_path) and os.path.exists(vectorizer_path):
        logger.info("Loading existing model...")
        model = pickle.load(open(model_path, 'rb'))
        vectorizer = pickle.load(open(vectorizer_path, 'rb'))
        return True
    
    logger.info("Creating new model...")
    
    # More comprehensive model for decoy generation
    model = RandomForestClassifier(n_estimators=100)
    
    # Extended training data for post-quantum algorithms
    sample_real = [
        "kyber768", "kyber1024", "saber", "ntru", "dilithium", 
        "falcon512", "falcon1024", "sphincs+", "picnic", "mceliece348864",
        "frodokem640", "frodokem976", "frodokem1344", "bike1", "bike2",
        "hqc128", "hqc192", "hqc256"
    ]
    
    # Generate fake data that looks similar to real data
    sample_fake = []
    for real in sample_real:
        # Generate 3 variations for each real algorithm
        for i in range(3):
            sample_fake.append(character_substitution_decoy(real, 0.7))
    
    # Combine real and fake data
    all_data = sample_real + sample_fake
    
    # Create features
    vectorizer = CountVectorizer(analyzer='char', ngram_range=(2, 5))
    X = vectorizer.fit_transform(all_data)
    
    # 1 for real algorithms, 0 for decoys
    y = [1] * len(sample_real) + [0] * len(sample_fake)
    
    # Train the model
    model.fit(X, y)
    
    # Save the model
    pickle.dump(model, open(model_path, 'wb'))
    pickle.dump(vectorizer, open(vectorizer_path, 'wb'))
    
    return True

def generate_decoy(target, complexity=5):
    """Generate a decoy based on the target string"""
    if not target:
        return ""
    
    # Complexity affects how different the decoy is (1-10)
    # Lower complexity means more similar to target
    similarity = 1.0 - (complexity / 20.0)  # Maps 1-10 to 0.95-0.5 similarity
    
    # Choose strategy based on complexity
    if complexity <= 3:
        # For low complexity, mostly character substitutions
        return character_substitution_decoy(target, similarity)
    elif complexity <= 7:
        # For medium complexity, use a mix of techniques
        if np.random.random() < 0.5:
            return character_substitution_decoy(target, similarity)
        else:
            return structural_variation_decoy(target, similarity)
    else:
        # For high complexity, use structural variations
        return structural_variation_decoy(target, similarity)

def character_substitution_decoy(target, similarity):
    """Generate decoy by substituting characters"""
    chars = list(target)
    
    # Determine how many characters to modify
    num_to_modify = int((1.0 - similarity) * len(target))
    num_to_modify = max(1, min(num_to_modify, len(target) - 1))  # At least 1, at most len-1
    
    # Choose positions to modify
    positions = np.random.choice(len(target), size=num_to_modify, replace=False)
    
    # Common substitutions that look similar
    substitutions = {
        'a': ['@', '4', 'á'],
        'b': ['6', '8', 'ß'],
        'e': ['3', 'é', 'ë'],
        'i': ['1', '!', 'í'],
        'l': ['1', '|', 'ł'],
        'o': ['0', 'ø', 'ó'],
        's': ['5', '$', 'š'],
        't': ['7', '+', 'ţ'],
        'z': ['2', 'ž', 'ż']
    }
    
    # For each position, modify the character
    for pos in positions:
        if target[pos].lower() in substitutions:
            # Use a similar-looking substitution
            options = substitutions[target[pos].lower()]
            chars[pos] = options[np.random.randint(0, len(options))]
        elif target[pos].isalpha():
            # Keep the case
            if target[pos].isupper():
                chars[pos] = chr(np.random.randint(65, 91))  # A-Z
            else:
                chars[pos] = chr(np.random.randint(97, 123))  # a-z
        elif target[pos].isdigit():
            # For digits, use a different digit
            options = [str(d) for d in range(10) if str(d) != target[pos]]
            chars[pos] = options[np.random.randint(0, len(options))]
        else:
            # For special characters, choose from a set of special chars
            special_chars = ['!', '@', '#', '$', '%', '^', '&', '*', '-', '_', '+', '=']
            chars[pos] = special_chars[np.random.randint(0, len(special_chars))]
    
    return ''.join(chars)

def structural_variation_decoy(target, similarity):
    """Generate decoy by changing the structure of the target"""
    # Patterns in post-quantum algorithm names
    patterns = [
        # kyber768 -> kyber-768
        (r'(\w+)(\d+)', r'\1-\2'),
        # kyber768 -> kyber_768
        (r'(\w+)(\d+)', r'\1_\2'),
        # ntru -> ntru-prime
        (r'(\w+)', r'\1-prime'),
        # dilithium -> dilithiumv2
        (r'(\w+)', r'\1v2'),
        # saber -> lightsaber
        (r'(\w+)', r'light\1'),
        # saber -> firesaber
        (r'(\w+)', r'fire\1'),
        # falcon512 -> falcon-512
        (r'(\w+)(\d+)', r'\1-\2'),
        # picnic -> picnic3
        (r'(\w+)', r'\1\d')
    ]
    
    # Choose a pattern at random
    pattern_idx = np.random.randint(0, len(patterns))
    pattern, replacement = patterns[pattern_idx]
    
    # Apply the pattern
    result = re.sub(pattern, replacement, target)
    
    # If result is the same as target, use character substitution as fallback
    if result == target:
        return character_substitution_decoy(target, similarity)
    
    return result

def calculate_similarity(decoy, target):
    """Calculate similarity between decoy and target"""
    return SequenceMatcher(None, decoy, target).ratio()

def evaluate_decoy(decoy, target):
    """Evaluate the effectiveness of a decoy"""
    # Calculate text similarity
    similarity = calculate_similarity(decoy, target)
    
    # Calculate probability of being real using the model
    if model is not None and vectorizer is not None:
        X = vectorizer.transform([decoy])
        real_probability = model.predict_proba(X)[0][1]
    else:
        real_probability = 0.5  # Default if model not available
    
    # Combine these metrics
    effectiveness = (similarity * 0.5) + (real_probability * 0.5)
    return effectiveness

def init_db_connection():
    """Initialize database connection"""
    db_path = os.environ.get('DB_PATH', '../database/pqcd.db')
    
    try:
        conn = sqlite3.connect(db_path)
        logger.info(f"Connected to database at {db_path}")
        return conn
    except sqlite3.Error as e:
        logger.error(f"Database connection failed: {e}")
        return None

@app.before_request
def before_request():
    """Initialize model before first request"""
    global model
    if model is None:
        load_or_create_model()

@app.route("/health", methods=["GET"])
def health():
    return jsonify({
        "status": "operational",
        "time": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
    })

@app.route("/generate", methods=["POST"])
def generate_decoys():
    start_time = time.time()
    
    data = request.get_json()
    if not data:
        return jsonify({"error": "No JSON data provided"}), 400

    target = data.get("target")
    if not target:
        return jsonify({"error": "Target string is required"}), 400

    complexity = int(data.get("complexity", 5))
    count = int(data.get("count", 10))
    
    # Validate inputs
    if not (1 <= complexity <= 10):
        complexity = 5
    
    if not (1 <= count <= 100):
        count = 10
    
    # Generate decoys
    decoys = []
    effectiveness_scores = []
    
    for _ in range(count):
        decoy = generate_decoy(target, complexity)
        decoys.append(decoy)
        
        # Calculate effectiveness score
        effectiveness = evaluate_decoy(decoy, target)
        effectiveness_scores.append(effectiveness)
    
    # Calculate generation time
    generation_time = int((time.time() - start_time) * 1000)  # in ms
    
    # Save to database if connection exists
    conn = init_db_connection()
    if conn:
        try:
            cursor = conn.cursor()
            for i, decoy in enumerate(decoys):
                cursor.execute(
                    "INSERT INTO decoys (decoy_text, target_text, complexity, effectiveness_score) VALUES (?, ?, ?, ?)",
                    (decoy, target, complexity, effectiveness_scores[i])
                )
            conn.commit()
            conn.close()
        except sqlite3.Error as e:
            logger.error(f"Database error: {e}")
    
    # Calculate average effectiveness
    avg_effectiveness = sum(effectiveness_scores) / len(effectiveness_scores) if effectiveness_scores else 0
    
    response = {
        "decoys": decoys,
        "target": target,
        "complexity": complexity,
        "count": count,
        "metrics": {
            "generation_time_ms": generation_time,
            "avg_effectiveness": round(avg_effectiveness, 2),
            "max_effectiveness": round(max(effectiveness_scores), 2) if effectiveness_scores else 0,
            "min_effectiveness": round(min(effectiveness_scores), 2) if effectiveness_scores else 0
        }
    }
    
    return jsonify(response)

@app.route("/evaluate", methods=["POST"])
def evaluate_decoys_endpoint():
    """Evaluate given decoys against a target"""
    data = request.get_json()
    if not data:
        return jsonify({"error": "No JSON data provided"}), 400

    target = data.get("target")
    decoys = data.get("decoys", [])
    
    if not target or not decoys:
        return jsonify({"error": "Target and decoys are required"}), 400
    
    # Evaluate each decoy
    results = []
    for decoy in decoys:
        effectiveness = evaluate_decoy(decoy, target)
        similarity = calculate_similarity(decoy, target)
        
        # Prepare result for this decoy
        result = {
            "decoy": decoy,
            "effectiveness": round(effectiveness, 2),
            "similarity": round(similarity, 2)
        }
        
        if model is not None and vectorizer is not None:
            X = vectorizer.transform([decoy])
            real_probability = model.predict_proba(X)[0][1]
            result["real_probability"] = round(real_probability, 2)
        
        results.append(result)
    
    # Sort by effectiveness
    results.sort(key=lambda x: x["effectiveness"], reverse=True)
    
    response = {
        "target": target,
        "evaluations": results,
        "metrics": {
            "avg_effectiveness": round(sum(r["effectiveness"] for r in results) / len(results), 2),
            "avg_similarity": round(sum(r["similarity"] for r in results) / len(results), 2)
        }
    }
    
    return jsonify(response)

@app.route('/generate', methods=['POST'])
def generate():
    data = request.json
    target = data.get('target', 'default')
    count = int(data.get('count', 5))
    complexity = int(data.get('complexity', 5))
    
    decoys = generate_decoys(target, complexity, count)
    
    return jsonify({
        "decoys": decoys,
        "metadata": {
            "model_version": "0.1.0",
            "complexity": complexity,
            "generation_time": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        }
    })

@app.route('/evaluate', methods=['POST'])
def evaluate():
    data = request.json
    target = data.get('target', '')
    decoys = data.get('decoys', [])
    
    results = []
    for decoy in decoys:
        similarity = random.uniform(0.5, 0.95)
        results.append({
            "decoy": decoy,
            "similarity": similarity,
            "effectiveness": similarity * (1 - similarity)  # Peak effectiveness around 0.5 similarity
        })
    
    return jsonify({
        "results": results,
        "metadata": {
            "model_version": "0.1.0",
            "evaluation_time": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        }
    })

def generate_decoys(target, complexity, count):
    """Generate decoys based on target text"""
    decoys = []
    
    for i in range(count):
        # Higher complexity = more different
        change_factor = complexity / 10.0
        
        if target.startswith('kyber'):
            decoy = f"kyber{random.choice(['512', '768', '1024', '1536'])}"
        elif target.startswith('saber'):
            decoy = random.choice(['lightsaber', 'saber', 'firesaber'])
        elif target.startswith('ntru'):
            params = random.choice(['hrss701', 'hps2048509', 'hps2048677', 'hps4096821'])
            decoy = f"ntru-{params}"
        elif target.startswith('dilithium'):
            mode = random.randint(1, 5)
            decoy = f"dilithium-mode{mode}"
        else:
            # Generate variation of the target
            chars = list(target)
            change_count = max(1, int(len(target) * change_factor * 0.5))
            for _ in range(change_count):
                pos = random.randint(0, len(chars) - 1)
                chars[pos] = random.choice(string.ascii_letters + string.digits)
            decoy = ''.join(chars)
            
        similarity = 1.0 - (change_factor * random.uniform(0.5, 1.0))
        similarity = max(0.1, min(0.9, similarity))  # Keep between 0.1 and 0.9
            
        decoys.append({
            "value": decoy,
            "similarity": similarity,
            "complexity": complexity
        })
    
    return decoys

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 5000))
    debug = os.environ.get("DEBUG", "False").lower() == "true"
    
    # Always ensure model is loaded
    load_or_create_model()
    
    app.run(host="127.0.0.1", port=port, debug=debug) 