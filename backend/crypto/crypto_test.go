package crypto

import (
	"bytes"
	"testing"
)

func TestKyberKeyGeneration(t *testing.T) {
	// Generate a Kyber key pair
	keyPair, err := GenerateKeyPair(AlgoKyber)
	if err != nil {
		t.Fatalf("Failed to generate Kyber key pair: %v", err)
	}

	// Verify the key pair
	if len(keyPair.PublicKey) != KyberPublicKeySize {
		t.Errorf("Expected public key size %d, got %d", KyberPublicKeySize, len(keyPair.PublicKey))
	}
	if len(keyPair.PrivateKey) != KyberPrivateKeySize {
		t.Errorf("Expected private key size %d, got %d", KyberPrivateKeySize, len(keyPair.PrivateKey))
	}
	if keyPair.Algorithm != AlgoKyber {
		t.Errorf("Expected algorithm %s, got %s", AlgoKyber, keyPair.Algorithm)
	}

	// Generate a fingerprint and verify it's not empty
	fingerprint := FingerPrint(keyPair.PublicKey)
	if fingerprint == "" {
		t.Error("Expected non-empty fingerprint")
	}
}

func TestEncryptionDecryption(t *testing.T) {
	// Generate a key pair
	keyPair, err := GenerateKeyPair(AlgoKyber)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Test data
	plaintext := []byte("This is a test message for post-quantum encryption")

	// Encrypt the data
	encrypted, err := Encrypt(plaintext, keyPair.PublicKey, keyPair.Algorithm)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Verify encrypted data
	if encrypted.Algorithm != keyPair.Algorithm {
		t.Errorf("Expected algorithm %s, got %s", keyPair.Algorithm, encrypted.Algorithm)
	}
	if len(encrypted.Nonce) != 24 {
		t.Errorf("Expected nonce length 24, got %d", len(encrypted.Nonce))
	}
	if len(encrypted.Ciphertext) <= KyberCiphertextSize {
		t.Errorf("Ciphertext too short, expected > %d bytes", KyberCiphertextSize)
	}

	// Decrypt the data
	decrypted, err := Decrypt(encrypted, keyPair.PrivateKey)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Verify decrypted data matches original plaintext
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted text does not match original plaintext")
		t.Errorf("Original: %s", plaintext)
		t.Errorf("Decrypted: %s", decrypted)
	}
}

func TestDecoyGeneration(t *testing.T) {
	// Generate a real key pair
	realKey, err := GenerateKeyPair(AlgoKyber)
	if err != nil {
		t.Fatalf("Failed to generate real key pair: %v", err)
	}

	// Generate decoy keys
	numDecoys := 3
	decoys, err := GenerateCognitiveDecoyKeys(realKey, numDecoys)
	if err != nil {
		t.Fatalf("Failed to generate decoy keys: %v", err)
	}

	// Verify the number of decoys
	if len(decoys) != numDecoys {
		t.Errorf("Expected %d decoys, got %d", numDecoys, len(decoys))
	}

	// Verify each decoy
	for i, decoy := range decoys {
		// Check algorithm matches
		if decoy.Algorithm != realKey.Algorithm {
			t.Errorf("Decoy %d: Expected algorithm %s, got %s", i, realKey.Algorithm, decoy.Algorithm)
		}

		// Check key sizes match
		if len(decoy.PublicKey) != len(realKey.PublicKey) {
			t.Errorf("Decoy %d: Expected public key size %d, got %d", i, len(realKey.PublicKey), len(decoy.PublicKey))
		}
		if len(decoy.PrivateKey) != len(realKey.PrivateKey) {
			t.Errorf("Decoy %d: Expected private key size %d, got %d", i, len(realKey.PrivateKey), len(decoy.PrivateKey))
		}

		// Check keys are different from real key
		if bytes.Equal(decoy.PublicKey, realKey.PublicKey) {
			t.Errorf("Decoy %d: Public key is identical to real key", i)
		}
		if bytes.Equal(decoy.PrivateKey, realKey.PrivateKey) {
			t.Errorf("Decoy %d: Private key is identical to real key", i)
		}

		// Check fingerprints are different
		realFingerprint := FingerPrint(realKey.PublicKey)
		decoyFingerprint := FingerPrint(decoy.PublicKey)
		if realFingerprint == decoyFingerprint {
			t.Errorf("Decoy %d: Fingerprint is identical to real key", i)
		}
	}
} 