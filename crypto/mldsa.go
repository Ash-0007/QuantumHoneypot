package crypto

import (
	"crypto/rand"
	"fmt"

	"github.com/cloudflare/circl/sign/dilithium/mode2"
)

// MLDSA65Provider implements the SignatureProvider interface for ML-DSA-65
type MLDSA65Provider struct{}

// NewMLDSA65Provider creates a new ML-DSA-65 provider
func NewMLDSA65Provider() *MLDSA65Provider {
	return &MLDSA65Provider{}
}

// Name returns the algorithm name
func (p *MLDSA65Provider) Name() Algorithm {
	return AlgMLDSA65
}

// KeyGen generates a new ML-DSA-65 key pair
func (p *MLDSA65Provider) KeyGen() (KeyPair, error) {
	// Generate key pair
	pk, sk, err := mode2.GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to generate ML-DSA-65 key pair: %w", err)
	}
	
	// Extract public and private keys as bytes
	publicKey := pk.Bytes()
	privateKey := sk.Bytes()

	return KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Algorithm:  AlgMLDSA65,
	}, nil
}

// Sign creates a signature for the given message using the private key
func (p *MLDSA65Provider) Sign(privateKeyBytes, message []byte) ([]byte, error) {
	// Parse private key from bytes
	sk := new(mode2.PrivateKey)
	if err := sk.UnmarshalBinary(privateKeyBytes); err != nil {
		return nil, fmt.Errorf("failed to parse ML-DSA-65 private key: %w", err)
	}
	
	// Sign the message
	signature, err := sk.Sign(rand.Reader, message, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message with ML-DSA-65: %w", err)
	}
	
	return signature, nil
}

// Verify checks if the signature is valid for the given message and public key
func (p *MLDSA65Provider) Verify(publicKeyBytes, message, signature []byte) (bool, error) {
	// Parse public key from bytes
	pk := new(mode2.PublicKey)
	if err := pk.UnmarshalBinary(publicKeyBytes); err != nil {
		return false, fmt.Errorf("failed to parse ML-DSA-65 public key: %w", err)
	}
	
	// Verify the signature
	valid := mode2.Verify(pk, message, signature)
	
	return valid, nil
} 