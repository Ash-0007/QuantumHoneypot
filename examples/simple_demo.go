package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// A simplified version of the crypto demo that doesn't rely on importing the crypto package
func main() {
	fmt.Println("=== Post-Quantum Cryptography Demo (Simplified) ===")
	fmt.Println("NOTE: This is a simulated implementation for demonstration purposes.")
	fmt.Println("      Do not use in production without a proper post-quantum library.")
	fmt.Println()

	// Test message to encrypt and decrypt
	message := []byte("Hello, post-quantum world!")
	fmt.Printf("Original message: %s\n\n", message)

	// Demo each supported algorithm
	algorithms := []string{"kyber", "saber", "ntru", "dilithium"}
	
	for _, alg := range algorithms {
		fmt.Printf("=== Testing %s algorithm ===\n", alg)
		
		// Generate key pair
		fmt.Printf("Generating %s key pair...\n", alg)
		publicKey, privateKey := generateSimulatedKeyPair(alg)
		
		// Get fingerprint
		fingerprint := generateFingerprint(publicKey)
		fmt.Printf("Key fingerprint: %s\n", fingerprint)
		fmt.Printf("Public key size: %d bytes\n", len(publicKey))
		fmt.Printf("Private key size: %d bytes\n", len(privateKey))
		
		// Encrypt message
		fmt.Println("Encrypting message...")
		ciphertext := simulateEncrypt(publicKey, message)
		
		// Decrypt message
		fmt.Println("Decrypting message...")
		decrypted := simulateDecrypt(publicKey, ciphertext)
		
		// Verify decryption
		success := string(decrypted) == string(message)
		if success {
			fmt.Printf("Decryption successful! Decrypted: %s\n", decrypted)
		} else {
			fmt.Printf("Decryption failed! Expected: %s, Got: %s\n", message, decrypted)
		}
		
		// Generate decoy keys
		fmt.Println("\nGenerating cognitive decoy keys...")
		decoyCount := 3
		decoyKeys := generateDecoyKeys(publicKey, decoyCount)
		
		fmt.Printf("Generated %d decoy keys\n", len(decoyKeys))
		
		// Test decryption with decoy keys (should fail)
		fmt.Println("Attempting decryption with decoy keys (should fail):")
		for i, decoyKey := range decoyKeys {
			// In a real implementation, this would fail because decoy keys can't decrypt
			// But in our simulated implementation, it will produce incorrect data
			decoyDecrypted := simulateDecrypt(decoyKey, ciphertext)
			
			match := string(decoyDecrypted) == string(message)
			if match {
				fmt.Printf("  Decoy #%d: WARNING - Decryption succeeded (this would indicate a security issue in a real implementation)\n", i+1)
			} else {
				fmt.Printf("  Decoy #%d: Decryption produced incorrect data as expected\n", i+1)
			}
		}
		
		fmt.Println()
	}
	
	fmt.Println("Demo completed!")
}

// Simulated key generation
func generateSimulatedKeyPair(algorithm string) ([]byte, []byte) {
	var publicKeySize, privateKeySize int
	
	switch algorithm {
	case "kyber":
		publicKeySize = 1184
		privateKeySize = 2400
	case "saber":
		publicKeySize = 992
		privateKeySize = 2304
	case "ntru":
		publicKeySize = 1138
		privateKeySize = 1450
	case "dilithium":
		publicKeySize = 1312
		privateKeySize = 2528
	default:
		publicKeySize = 1184
		privateKeySize = 2400
	}
	
	// Generate a seed for deterministic key generation
	seed := make([]byte, 32)
	rand.Read(seed)
	
	// Generate keys deterministically from seed
	publicKey := deriveKey(seed, publicKeySize)
	privateKey := deriveKey(append(seed, 0x01), privateKeySize)
	
	return publicKey, privateKey
}

// Derive a key of specified size from a seed
func deriveKey(seed []byte, size int) []byte {
	key := make([]byte, size)
	
	// Use the seed to fill the key
	for i := 0; i < size; i += 32 {
		h := sha256.New()
		h.Write(seed)
		h.Write([]byte{byte(i / 32)})
		chunk := h.Sum(nil)
		
		copySize := 32
		if i+32 > size {
			copySize = size - i
		}
		
		copy(key[i:i+copySize], chunk[:copySize])
	}
	
	return key
}

// Generate a fingerprint for a key
func generateFingerprint(key []byte) string {
	hash := sha256.Sum256(key)
	return fmt.Sprintf("%x", hash[:8])
}

// Simulate encryption
func simulateEncrypt(publicKey, plaintext []byte) []byte {
	// Generate a shared secret from the public key
	sharedSecret := sha256.Sum256(publicKey)
	
	// XOR the plaintext with the derived key
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		ciphertext[i] = plaintext[i] ^ sharedSecret[i%len(sharedSecret)]
	}
	
	return ciphertext
}

// Simulate decryption
func simulateDecrypt(publicKey, ciphertext []byte) []byte {
	// In a real implementation, we would use the private key to derive the same shared secret
	// Here we just use the public key to get the same shared secret as in encryption
	sharedSecret := sha256.Sum256(publicKey)
	
	// XOR the ciphertext with the derived key
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i++ {
		plaintext[i] = ciphertext[i] ^ sharedSecret[i%len(sharedSecret)]
	}
	
	return plaintext
}

// Generate decoy keys
func generateDecoyKeys(publicKey []byte, count int) [][]byte {
	decoys := make([][]byte, count)
	
	for i := 0; i < count; i++ {
		decoy := make([]byte, len(publicKey))
		copy(decoy, publicKey)
		
		// Make some small changes to create a decoy
		for j := 0; j < 10; j++ {
			pos := (i * j) % len(decoy)
			decoy[pos] = decoy[pos] ^ 0xFF
		}
		
		decoys[i] = decoy
	}
	
	return decoys
} 