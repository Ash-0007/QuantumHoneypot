# Post-Quantum Cryptography Examples

This directory contains example code demonstrating how to use the post-quantum cryptography implementation in this project.

## Available Examples

### `crypto_demo.go`

A comprehensive demo of the post-quantum cryptography implementation, showing:

- Key generation for different algorithms (Kyber, Saber, NTRU, Dilithium)
- Encryption and decryption of messages
- Cognitive decoy generation and testing

To run this example:

```bash
# From the project root directory
cd pqcd/examples
go run crypto_demo.go
```

### `simple_demo.go`

A simplified, self-contained demo that doesn't require importing the crypto package. This demo:

- Simulates key generation for different post-quantum algorithms with realistic key sizes
- Implements a simple encryption/decryption mechanism using SHA-256 and XOR
- Demonstrates cognitive decoy key generation and testing

To run this example:

```bash
# From the project root directory
cd pqcd/examples
go run simple_demo.go
```

## Implementation Details

The current implementation is a functional simulation of post-quantum cryptography algorithms for demonstration purposes. While we initially attempted to use the CloudFlare CIRCL library, we encountered compatibility issues with the specific version available in this environment.

Key features include:

- **Simulated Key Generation**: Generates key pairs with realistic sizes matching actual post-quantum algorithms
- **Simulated Encryption/Decryption**: Uses SHA-256 for key derivation and XOR-based encryption
- **Cognitive Decoy Generation**: Creates decoy keys that appear similar to real keys
- **Fingerprinting**: Generates unique fingerprints for keys using SHA-256 hashing

For production use, this implementation should be replaced with actual post-quantum cryptography libraries.

For more information, see the main project README. 