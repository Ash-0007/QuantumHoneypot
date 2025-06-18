package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"fmt"
)

// ECDHProvider implements the KEMProvider interface for ECDH
type ECDHProvider struct{}

// NewECDHProvider creates a new ECDH provider
func NewECDHProvider() *ECDHProvider {
	return &ECDHProvider{}
}

// Name returns the algorithm name
func (p *ECDHProvider) Name() Algorithm {
	return AlgECDH
}

// KeyGen generates a new ECDH key pair
func (p *ECDHProvider) KeyGen() (KeyPair, error) {
	// Generate private key using P-256 curve
	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to generate ECDH key pair: %w", err)
	}
	
	// Extract public and private keys as bytes
	publicKeyBytes := privateKey.PublicKey().Bytes()
	
	// Marshal private key to PKCS8 format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to marshal ECDH private key: %w", err)
	}
	
	return KeyPair{
		PublicKey:  publicKeyBytes,
		PrivateKey: privateKeyBytes,
		Algorithm:  AlgECDH,
	}, nil
}

// Encapsulate generates an ephemeral key pair, computes shared secret, and returns the ephemeral public key as ciphertext
func (p *ECDHProvider) Encapsulate(publicKeyBytes []byte) ([]byte, []byte, error) {
	// Parse recipient's public key
	recipientPubKey, err := ecdh.P256().NewPublicKey(publicKeyBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse ECDH public key: %w", err)
	}
	
	// Generate ephemeral key pair
	ephemeralKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate ephemeral ECDH key: %w", err)
	}
	
	// Compute shared secret
	sharedSecret, err := ephemeralKey.ECDH(recipientPubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute ECDH shared secret: %w", err)
	}
	
	// Return ephemeral public key as ciphertext
	ephemeralPubKeyBytes := ephemeralKey.PublicKey().Bytes()
	
	return ephemeralPubKeyBytes, sharedSecret, nil
}

// Decapsulate computes the shared secret using the private key and the ephemeral public key (ciphertext)
func (p *ECDHProvider) Decapsulate(privateKeyBytes, ciphertextBytes []byte) ([]byte, error) {
	// Parse private key from PKCS8 format
	privKeyInterface, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ECDH private key: %w", err)
	}
	
	// Cast to ECDH private key
	privateKey, ok := privKeyInterface.(*ecdh.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid ECDH private key type")
	}
	
	// Parse ephemeral public key from ciphertext
	ephemeralPubKey, err := ecdh.P256().NewPublicKey(ciphertextBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ephemeral ECDH public key: %w", err)
	}
	
	// Compute shared secret
	sharedSecret, err := privateKey.ECDH(ephemeralPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute ECDH shared secret: %w", err)
	}
	
	return sharedSecret, nil
} 