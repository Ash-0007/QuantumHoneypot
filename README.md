# Post-Quantum Cryptography Evaluation Framework

This project implements a framework for evaluating Post-Quantum Cryptography (PQC) algorithms within a RESTful API environment, as described in the research paper "A Framework for Evaluating Post-Quantum Cryptography in Go-Based Secure APIs".

## Features

- Implementation of NIST-standardized PQC algorithms:
  - ML-KEM-768 (based on Kyber768) for key encapsulation
  - ML-DSA-65 (based on Dilithium2) for digital signatures
- Implementation of classical counterparts for comparison:
  - ECDH with P-256 curve
  - ECDSA with P-256 curve
- RESTful API endpoints for all cryptographic operations
- Performance measurement and benchmarking
- AI-driven adaptive threat intelligence and honeypot system

## Requirements

- Go 1.21 or newer
- Dependencies (automatically installed via Go modules):
  - github.com/cloudflare/circl
  - github.com/gorilla/mux
  - github.com/sirupsen/logrus
  - golang.org/x/time

## Installation

```bash
git clone https://github.com/compsec/pqcd.git
cd pqcd
go build
```

## Usage

### Starting the Server

```bash
# Start with default settings (port 8080, AI security disabled)
./pqcd

# Start with AI security enabled
./pqcd --enable-ai

# Change port and log level
./pqcd --port 9000 --log-level debug
```

### API Endpoints

#### Key Encapsulation (ML-KEM-768 and ECDH)

**Generate Key Pair:**
```
POST /api/{alg}/keygen
```

**Encapsulate (Generate Shared Secret):**
```
POST /api/{alg}/encapsulate
{
  "publicKey": "hex-encoded-public-key"
}
```

**Decapsulate (Recover Shared Secret):**
```
POST /api/{alg}/decapsulate
{
  "privateKey": "hex-encoded-private-key",
  "ciphertext": "hex-encoded-ciphertext"
}
```

Where `{alg}` is one of:
- `ml-kem-768` (post-quantum)
- `ecdh` (classical)

#### Digital Signatures (ML-DSA-65 and ECDSA)

**Generate Key Pair:**
```
POST /api/{alg}/keygen
```

**Sign:**
```
POST /api/{alg}/sign
{
  "privateKey": "hex-encoded-private-key",
  "message": "message-to-sign"
}
```

**Verify:**
```
POST /api/{alg}/verify
{
  "publicKey": "hex-encoded-public-key",
  "message": "message-that-was-signed",
  "signature": "hex-encoded-signature"
}
```

Where `{alg}` is one of:
- `ml-dsa-65` (post-quantum)
- `ecdsa` (classical)

### Metrics

View performance metrics:
```
GET /api/metrics
```

## AI Security Layer

The AI-driven security layer includes:

1. **Feature Engineering**: Extracts features from API requests
2. **Anomaly Detection**: Uses a statistical model (in production, a deep autoencoder) to detect anomalous behavior
3. **Threat Classification**: Classifies anomalies into threat types
4. **Dynamic Response Engine**: Decides on appropriate actions:
   - Pass: Allow the request to proceed normally
   - Throttle: Introduce artificial delays to slow attackers
   - Deceive: Return validly formatted but cryptographically incorrect responses
   - Redirect: Transparently redirect to a honeypot system for threat intelligence gathering

## Paper Citation

If you use this framework in your research, please cite our paper:
```
@inproceedings{doe2025framework,
  title={A Framework for Evaluating Post-Quantum Cryptography in Go-Based Secure APIs},
  author={Doe, John and Smith, Jane and Brown, Robert},
  booktitle={Proceedings of the International Conference on Cryptographic Engineering},
  year={2025},
  organization={IEEE}
}
```

## License

MIT 