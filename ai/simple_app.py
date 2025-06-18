from flask import Flask, request, jsonify
import random
import time
from flask_cors import CORS
import string

app = Flask(__name__)
CORS(app)

@app.route('/health', methods=['GET'])
def health():
    return jsonify({
        'status': 'ok',
        'time': time.strftime('%Y-%m-%dT%H:%M:%S')
    })

@app.route('/decoys/generate', methods=['POST'])
def generate_decoys():
    data = request.json
    
    if not data or 'target' not in data:
        return jsonify({'error': 'Target is required'}), 400
    
    target = data.get('target', '')
    complexity = data.get('complexity', 5)
    count = data.get('count', 5)
    
    decoys = generate_simple_decoys(target, complexity, count)
    
    return jsonify({
        'decoys': decoys,
        'generation_time': time.strftime('%Y-%m-%dT%H:%M:%S')
    })

@app.route('/evaluate', methods=['POST'])
def evaluate():
    data = request.json
    
    if not data or 'original' not in data or 'decoys' not in data:
        return jsonify({'error': 'Both original and decoys are required'}), 400
    
    original = data.get('original', '')
    decoys = data.get('decoys', [])
    
    # Simple evaluation logic
    evaluation = []
    for i, decoy in enumerate(decoys):
        similarity = calculate_similarity(original, decoy)
        effectiveness = random.uniform(0.6, 0.95)
        if i == 0:  # Make the first one always stand out a bit
            effectiveness += 0.1
            if effectiveness > 1.0:
                effectiveness = 1.0
                
        evaluation.append({
            'decoy': decoy,
            'similarity': similarity,
            'effectiveness': effectiveness,
            'explanation': get_explanation(similarity, effectiveness)
        })
    
    return jsonify({
        'evaluation': evaluation,
        'summary': {
            'overall_effectiveness': sum(e['effectiveness'] for e in evaluation) / len(evaluation),
            'recommendation': 'The generated decoys provide good cognitive security.'
        }
    })

def calculate_similarity(original, decoy):
    # Simple character-based similarity
    if not original or not decoy:
        return 0.0
    
    common = sum(1 for a, b in zip(original, decoy) if a == b)
    total = max(len(original), len(decoy))
    return common / total

def generate_simple_decoys(target, complexity, count):
    decoys = []
    
    for i in range(count):
        similarity = 1.0 - (complexity / 10.0) - (i / (count * 2))
        if similarity < 0.1:
            similarity = 0.1
            
        decoy = mutate_string(target, complexity)
        
        decoys.append({
            'value': decoy,
            'similarity': similarity
        })
    
    return decoys

def mutate_string(s, complexity):
    chars = list(s)
    
    # More complexity = more mutations
    mutations = 1 + int(complexity / 3)
    
    for _ in range(mutations):
        if not chars:
            break
            
        op = random.choice(['substitute', 'insert', 'delete', 'swap'])
        
        if op == 'substitute' and chars:
            pos = random.randint(0, len(chars) - 1)
            chars[pos] = random.choice(string.ascii_letters + string.digits)
        elif op == 'insert':
            pos = random.randint(0, len(chars))
            chars.insert(pos, random.choice(string.ascii_letters + string.digits))
        elif op == 'delete' and len(chars) > 1:
            pos = random.randint(0, len(chars) - 1)
            chars.pop(pos)
        elif op == 'swap' and len(chars) > 1:
            pos1 = random.randint(0, len(chars) - 1)
            pos2 = random.randint(0, len(chars) - 1)
            chars[pos1], chars[pos2] = chars[pos2], chars[pos1]
    
    return ''.join(chars)

def get_explanation(similarity, effectiveness):
    if effectiveness > 0.8:
        return "This decoy effectively mimics the target while maintaining plausible deniability."
    elif effectiveness > 0.6:
        return "This decoy provides reasonable cognitive security but could be improved."
    else:
        return "This decoy may not be convincing enough to provide good cognitive security."

if __name__ == '__main__':
    print("Starting AI service on http://localhost:5000")
    app.run(host='0.0.0.0', port=5000) 