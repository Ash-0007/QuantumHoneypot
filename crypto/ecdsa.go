package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// ECDSAProvider implements the SignatureProvider interface for ECDSA
type ECDSAProvider struct{}

// NewECDSAProvider creates a new ECDSA provider
func NewECDSAProvider() *ECDSAProvider {
	return &ECDSAProvider{}
}

// Name returns the algorithm name
func (p *ECDSAProvider) Name() Algorithm {
	return AlgECDSA
}

// KeyGen generates a new ECDSA key pair
func (p *ECDSAProvider) KeyGen() (KeyPair, error) {
	// Generate key pair using P-256 curve
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to generate ECDSA key pair: %w", err)
	}
	
	// Extract public key as bytes (X and Y coordinates)
	publicKeyBytes := elliptic.MarshalCompressed(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)
	
	// Extract private key as bytes (just the D value)
	privateKeyBytes := privateKey.D.Bytes()
	
	return KeyPair{
		PublicKey:  publicKeyBytes,
		PrivateKey: privateKeyBytes,
		Algorithm:  AlgECDSA,
	}, nil
}

// Sign creates a signature for the given message using the private key
func (p *ECDSAProvider) Sign(privateKeyBytes, message []byte) ([]byte, error) {
	// Parse private key
	d := new(big.Int).SetBytes(privateKeyBytes)
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: d,
	}
	
	// Compute public key from private key
	privateKey.PublicKey.X, privateKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(privateKeyBytes)
	
	// Hash the message
	digest := sha256.Sum256(message)
	
	// Sign the digest
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, digest[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign message with ECDSA: %w", err)
	}
	
	// Encode signature as bytes (R || S)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	signature := make([]byte, 64)
	
	// Ensure R and S values are padded to 32 bytes
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):64], sBytes)
	
	return signature, nil
}

// Verify checks if the signature is valid for the given message and public key
func (p *ECDSAProvider) Verify(publicKeyBytes, message, signature []byte) (bool, error) {
	if len(signature) != 64 {
		return false, fmt.Errorf("invalid ECDSA signature length: expected 64 bytes, got %d", len(signature))
	}
	
	// Parse public key
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKeyBytes)
	if x == nil {
		return false, fmt.Errorf("failed to unmarshal ECDSA public key")
	}
	
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	
	// Parse signature
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	
	// Hash the message
	digest := sha256.Sum256(message)
	
	// Verify the signature
	valid := ecdsa.Verify(publicKey, digest[:], r, s)
	
	return valid, nil
} 