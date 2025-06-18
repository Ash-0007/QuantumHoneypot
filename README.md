# Post-Quantum Cognitive Decoys (PQCD)

A 24-hour sprint implementation of a post-quantum security system leveraging cognitive decoys to enhance protection against quantum computing threats.

## Features

- **Post-Quantum Cryptography**: Simulate post-quantum cryptographic algorithms (Kyber, Saber, NTRU)
- **Cognitive Decoy Generation**: AI-powered algorithms generate convincing decoys to confuse attackers
- **Key Management**: Generate and manage post-quantum key pairs with decoy keys
- **Encryption/Decryption**: Secure message encryption and decryption
- **Interactive UI**: React-based interface for all cryptographic operations

## Project Structure

- `backend/`: Go backend server with cryptographic functions
  - `main.go`: HTTP server with REST API endpoints
  - `simple.go`: Simplified backend implementation for testing
  - `crypto/`: Post-quantum cryptography simulation
- `ai/`: Python ML models for cognitive decoy generation
  - `app.py`: Flask API for decoy generation
  - `db.py`: Database connection utilities
  - `models/`: Trained ML models for decoy evaluation
- `frontend/`: React-based demo interface
  - `src/App.js`: Main React application
  - `public/`: Static assets
- `database/`: SQLite database schemas and migrations
  - `schema.sql`: Database schema definition
  - `init_db.py`: Database initialization script
- `deployment/`: Docker compose files for one-command deployment
  - `docker-compose.yml`: Service orchestration configuration
  - `Dockerfile.*`: Container definitions for each service
- `docs/`: Project documentation

## Quick Start

### Using Scripts

```bash
# Clone the repository
git clone https://github.com/yourusername/pqcd.git
cd pqcd

# Start all services (Linux/macOS)
chmod +x run.sh
./run.sh

# Start all services (Windows)
.\run.ps1
```

After startup, access the application:

- Go Backend API: http://localhost:8083
- AI Service API: http://localhost:5000
- Test Frontend: Open test_frontend.html in your browser

### Manual Setup

#### Backend (Go)

```bash
cd pqcd/backend
go mod download github.com/mattn/go-sqlite3
go mod download github.com/rs/cors
go run main.go
```

#### AI Service (Python)

```bash
cd pqcd/ai
pip install -r requirements.txt
python app.py
```

#### Database Initialization

```bash
cd pqcd/database
python init_db.py
```

## API Endpoints

### Go Backend (port 8083)

- `GET /api/health`: Health check endpoint
- `GET /api/status`: Service status information
- `POST /api/decoys/generate`: Generate cognitive decoys
- `POST /api/keys/generate`: Generate post-quantum key pairs
- `POST /api/encrypt`: Encrypt a message
- `POST /api/decrypt`: Decrypt a message

### AI Service (port 5000)

- `GET /health`: Health check endpoint
- `POST /generate`: Generate cognitive decoys
- `POST /evaluate`: Evaluate decoy effectiveness

## Development

### Prerequisites

- Go 1.19+
- Python 3.9+
- SQLite

### Testing the API

You can use the included test_frontend.html file to test the API endpoints, or use curl:

```bash
# Test backend health
curl http://localhost:8083/api/health

# Generate keys
curl -X POST http://localhost:8083/api/keys/generate -H "Content-Type: application/json" -d '{"algorithm":"kyber","count":5}'

# Generate decoys
curl -X POST http://localhost:8083/api/decoys/generate -H "Content-Type: application/json" -d '{"target":"kyber768","complexity":5,"count":3}'

# Test AI service health
curl http://localhost:5000/health
```

## Next Steps for Production Implementation

1. **Replace Mock Implementations**
   - Integrate actual post-quantum cryptography libraries
   - Implement real key generation algorithms (Kyber, Saber, NTRU)
   - Enhance the AI model with more sophisticated decoy generation techniques

2. **Security Enhancements**
   - Add user authentication system (JWT or API keys)
   - Implement secure key storage with encryption
   - Enable TLS for all communications
   - Add rate limiting and other security measures

3. **Performance Optimization**
   - Optimize database queries for better performance
   - Implement caching for frequently accessed data
   - Add connection pooling for database access

4. **Scalability Improvements**
   - Implement horizontal scaling for all components
   - Add load balancing for high availability
   - Create stateless services for better resilience

5. **Monitoring and Observability**
   - Add comprehensive logging system
   - Implement metrics collection and dashboards
   - Create alerting system for service issues
   - Add performance monitoring and tracing

## About Post-Quantum Cognitive Decoys

Post-Quantum Cognitive Decoys (PQCD) is a novel security approach that combines post-quantum cryptography with cognitive science principles. By generating convincing decoys alongside real cryptographic keys, the system increases the computational cost and cognitive load for attackers, even those with access to quantum computing resources.

## Post-Quantum Cryptography Implementation

The current implementation includes a functional simulation of post-quantum cryptography algorithms for demonstration purposes. While we initially attempted to use the CloudFlare CIRCL library, we encountered compatibility issues with the specific version available in this environment.

Key features of our implementation include:

- **Simulated Key Generation**: Generates key pairs with realistic sizes matching actual post-quantum algorithms:
  - Kyber (1184 bytes public key, 2400 bytes private key)
  - Saber (992 bytes public key, 2304 bytes private key)
  - NTRU (699 bytes public key, 935 bytes private key)
  - Dilithium (1312 bytes public key, 2528 bytes private key)

- **Simulated Encryption/Decryption**: Implements a secure encryption/decryption mechanism using SHA-256 for key derivation and XOR-based encryption, with proper nonce handling.

- **Cognitive Decoy Generation**: Creates decoy keys that appear similar to real keys but have subtle differences, making it difficult for attackers to identify the genuine keys.

- **Fingerprinting**: Generates unique fingerprints for keys using SHA-256 hashing to enable key identification without revealing the entire key.

For production use, this implementation should be replaced with actual post-quantum cryptography libraries such as:
- CloudFlare's CIRCL library (with proper version compatibility)
- Open Quantum Safe (liboqs)
- NIST PQC standardized algorithm implementations
- BoringSSL with PQC support

## License

MIT License 