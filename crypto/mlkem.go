package crypto

import (
	"fmt"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/schemes"
)

// MLKEM768Provider implements the KEMProvider interface for ML-KEM-768
type MLKEM768Provider struct {
	scheme kem.Scheme
}

// NewMLKEM768Provider creates a new ML-KEM-768 provider
func NewMLKEM768Provider() *MLKEM768Provider {
	return &MLKEM768Provider{
		scheme: schemes.ByName("Kyber768"),
	}
}

// Name returns the algorithm name
func (p *MLKEM768Provider) Name() Algorithm {
	return AlgMLKEM768
}

// KeyGen generates a new ML-KEM-768 key pair
func (p *MLKEM768Provider) KeyGen() (KeyPair, error) {
	// Generate key pair
	pk, sk, err := p.scheme.GenerateKeyPair()
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to generate ML-KEM-768 key pair: %w", err)
	}

	// Extract public and private keys as bytes
	publicKey, err := pk.MarshalBinary()
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to marshal public key: %w", err)
	}

	privateKey, err := sk.MarshalBinary()
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to marshal private key: %w", err)
	}

	return KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Algorithm:  AlgMLKEM768,
	}, nil
}

// Encapsulate generates a shared secret and ciphertext using the recipient's public key
func (p *MLKEM768Provider) Encapsulate(publicKeyBytes []byte) ([]byte, []byte, error) {
	// Parse public key from bytes
	pk, err := p.scheme.UnmarshalBinaryPublicKey(publicKeyBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse ML-KEM-768 public key: %w", err)
	}

	// Encapsulate to generate ciphertext and shared secret
	ct, ss, err := p.scheme.Encapsulate(pk)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encapsulate using ML-KEM-768: %w", err)
	}

	return ct, ss, nil
}

// Decapsulate recovers the shared secret from the ciphertext using the private key
func (p *MLKEM768Provider) Decapsulate(privateKeyBytes, ciphertextBytes []byte) ([]byte, error) {
	// Parse private key from bytes
	sk, err := p.scheme.UnmarshalBinaryPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ML-KEM-768 private key: %w", err)
	}

	// Decapsulate to recover the shared secret
	ss, err := p.scheme.Decapsulate(sk, ciphertextBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decapsulate using ML-KEM-768: %w", err)
	}

	return ss, nil
} 