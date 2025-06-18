package crypto

import "time"

// Algorithm represents a cryptographic algorithm type
type Algorithm string

const (
	// Key Encapsulation Mechanisms
	AlgMLKEM768 Algorithm = "ml-kem-768"
	AlgECDH     Algorithm = "ecdh"

	// Digital Signature Algorithms
	AlgMLDSA65 Algorithm = "ml-dsa-65"
	AlgECDSA   Algorithm = "ecdsa"
)

// CryptoProvider is an interface for crypto operations
type CryptoProvider interface {
	// Name returns the name of the algorithm
	Name() Algorithm
	
	// KeyGen generates a new key pair
	KeyGen() (KeyPair, error)
}

// KeyPair represents a generic key pair
type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
	Algorithm  Algorithm
}

// OperationResult contains details about a cryptographic operation
type OperationResult struct {
	Algorithm Algorithm
	Duration  time.Duration
	Success   bool
	KeySize   int
	DataSize  int // Size of either signature or ciphertext
}

// KEMProvider is an interface for KEM operations
type KEMProvider interface {
	CryptoProvider
	
	// Encapsulate generates a shared secret for the given public key
	Encapsulate(publicKey []byte) (ciphertext []byte, sharedSecret []byte, err error)
	
	// Decapsulate recovers the shared secret from the ciphertext using the private key
	Decapsulate(privateKey, ciphertext []byte) (sharedSecret []byte, err error)
}

// SignatureProvider is an interface for digital signature operations
type SignatureProvider interface {
	CryptoProvider
	
	// Sign creates a signature for the given message using the private key
	Sign(privateKey, message []byte) (signature []byte, err error)
	
	// Verify checks if the signature is valid for the given message and public key
	Verify(publicKey, message, signature []byte) (valid bool, err error)
} 