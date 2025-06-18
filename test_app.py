from flask import Flask, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app)  # Enable CORS for all routes

@app.route("/health", methods=["GET"])
def health():
    return jsonify({"status": "ok"})

@app.route("/api/test", methods=["GET"])
def test():
    return jsonify({
        "message": "Post-Quantum Cognitive Decoys API is working",
        "decoys": ["test-decoy-1", "test-decoy-2", "test-decoy-3"]
    })

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=True) 