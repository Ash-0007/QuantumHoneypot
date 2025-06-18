# Post-Quantum Cognitive Decoys - Architecture

This document describes the architecture of the Post-Quantum Cognitive Decoys (PQCD) system.

## System Overview

The PQCD system consists of three main components:

1. **Go Backend**: Handles cryptographic operations and API endpoints
2. **Python AI Service**: Generates cognitive decoys using machine learning
3. **React Frontend**: Provides a user interface for interacting with the system

These components communicate with each other through HTTP APIs and share a common SQLite database.

## Architecture Diagram

```
┌─────────────────┐     HTTP     ┌─────────────────┐
│                 │◄────────────►│                 │
│  React Frontend │              │   Go Backend    │
│  (port: 3000)   │              │   (port: 8082)  │
│                 │              │                 │
└─────────────────┘              └────────┬────────┘
                                          │
                                          │ HTTP
                                          │
                                 ┌────────▼────────┐
                                 │                 │
                                 │  AI Service     │
                                 │  (port: 5000)   │
                                 │                 │
                                 └────────┬────────┘
                                          │
                                          │
                                 ┌────────▼────────┐
                                 │                 │
                                 │  SQLite DB      │
                                 │  (pqcd.db)      │
                                 │                 │
                                 └─────────────────┘
```

## Component Details

### Go Backend

The Go backend is responsible for:

- Handling API requests
- Performing cryptographic operations
- Managing database interactions
- Communicating with the AI service

**Key Files:**
- `main.go`: Main application entry point and API handlers
- `crypto/pq.go`: Post-quantum cryptography simulation

**API Endpoints:**
- `/api/health`: Health check endpoint
- `/api/status`: Service status information
- `/api/keys/generate`: Generate post-quantum key pairs
- `/api/decoys/generate`: Generate cognitive decoys
- `/api/encrypt`: Encrypt a message
- `/api/decrypt`: Decrypt a message

### Python AI Service

The AI service is responsible for:

- Generating cognitive decoys using machine learning
- Evaluating the effectiveness of decoys
- Training and updating decoy generation models

**Key Files:**
- `app.py`: Flask application with API endpoints
- `db.py`: Database connection module
- `models/`: Directory containing trained ML models

**API Endpoints:**
- `/health`: Health check endpoint
- `/generate`: Generate cognitive decoys
- `/evaluate`: Evaluate decoy effectiveness

### React Frontend

The frontend provides a user interface for:

- Generating post-quantum key pairs
- Creating and managing cognitive decoys
- Encrypting and decrypting messages
- Visualizing decoy effectiveness

**Key Files:**
- `src/App.js`: Main React application
- `src/components/`: React components for different features

### SQLite Database

The database stores:

- Generated key pairs (real and decoys)
- Cognitive decoys
- Event logs
- User information

**Key Tables:**
- `key_pairs`: Stores generated key pairs
- `decoys`: Stores generated decoys
- `event_logs`: Stores system events
- `users`: Stores user information

## Data Flow

### Key Generation Flow

1. User requests key generation via the frontend
2. Frontend sends request to Go backend
3. Backend generates a post-quantum key pair
4. Backend stores the key pair in the database
5. Backend generates decoy keys using the crypto module
6. Backend stores decoy keys in the database
7. Backend returns the public key and fingerprint to the frontend

### Decoy Generation Flow

1. User requests decoy generation via the frontend
2. Frontend sends request to Go backend
3. Backend forwards request to AI service
4. AI service generates decoys using machine learning
5. AI service returns decoys to the backend
6. Backend stores decoys in the database
7. Backend returns decoys to the frontend

### Encryption Flow

1. User inputs a message and selects a public key
2. Frontend sends encryption request to the backend
3. Backend encrypts the message using the selected public key
4. Backend logs the encryption event
5. Backend returns the encrypted data to the frontend

### Decryption Flow

1. User inputs ciphertext, private key, and nonce
2. Frontend sends decryption request to the backend
3. Backend decrypts the message using the private key
4. Backend logs the decryption event
5. Backend returns the plaintext to the frontend

## Deployment Architecture

The system can be deployed using Docker Compose, which creates containers for each component:

```
┌─────────────────────────────────────────────────────────┐
│                      Docker Host                        │
│                                                         │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│   │             │    │             │    │             │ │
│   │  Frontend   │    │   Backend   │    │  AI Service │ │
│   │  Container  │    │  Container  │    │  Container  │ │
│   │             │    │             │    │             │ │
│   └─────────────┘    └─────────────┘    └─────────────┘ │
│                                                         │
│                     ┌─────────────┐                     │
│                     │             │                     │
│                     │    DB Init  │                     │
│                     │  Container  │                     │
│                     │             │                     │
│                     └─────────────┘                     │
│                                                         │
│   ┌─────────────────────────────────────────────────┐   │
│   │                 Shared Volume                    │   │
│   │                   (pqcd.db)                      │   │
│   └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Security Considerations

The current implementation is a proof of concept and has several security limitations:

1. **No Authentication**: The API endpoints do not implement authentication
2. **Simulated Cryptography**: The post-quantum algorithms are simulated rather than using actual implementations
3. **Plaintext Storage**: Keys are stored in plaintext in the database
4. **No TLS**: Communications between components are not encrypted

For a production deployment, the following improvements would be necessary:

1. Add authentication using JWT or API keys
2. Implement actual post-quantum cryptography libraries
3. Use secure key storage with proper encryption
4. Enable TLS for all communications
5. Add rate limiting and other security measures

## Extensibility

The system is designed to be extensible in several ways:

1. **Additional Algorithms**: New post-quantum algorithms can be added to the crypto module
2. **Enhanced Decoy Generation**: The AI service can be extended with more sophisticated decoy generation algorithms
3. **User Management**: A user management system can be added for multi-user support
4. **Monitoring**: Monitoring and alerting can be added for operational visibility
5. **High Availability**: The components can be scaled horizontally for high availability 