package main

import (
	"fmt"
	"github.com/pqcd/backend/crypto"
)

func main() {
	fmt.Println("=== Post-Quantum Cryptography Demo ===")
	fmt.Println("NOTE: This is a simulated implementation for demonstration purposes.")
	fmt.Println("      Do not use in production without a proper post-quantum library.")
	fmt.Println()

	// Test message to encrypt and decrypt
	message := []byte("Hello, post-quantum world!")
	fmt.Printf("Original message: %s\n\n", message)

	// Demo each supported algorithm
	algorithms := []string{
		crypto.AlgoKyber,
		crypto.AlgoSaber,
		crypto.AlgoNTRU,
		crypto.AlgoDilithium,
	}
	
	for _, alg := range algorithms {
		fmt.Printf("=== Testing %s algorithm ===\n", alg)
		
		// Generate key pair
		fmt.Printf("Generating %s key pair...\n", alg)
		keyPair, err := crypto.GenerateKeyPair(alg)
		if err != nil {
			fmt.Printf("Error generating key pair: %v\n", err)
			continue
		}
		
		// Get fingerprint
		fingerprint := crypto.FingerPrint(keyPair.PublicKey)
		fmt.Printf("Key fingerprint: %s\n", fingerprint)
		fmt.Printf("Public key size: %d bytes\n", len(keyPair.PublicKey))
		fmt.Printf("Private key size: %d bytes\n", len(keyPair.PrivateKey))
		
		// Encrypt message
		fmt.Println("Encrypting message...")
		encrypted, err := crypto.Encrypt(message, keyPair.PublicKey, keyPair.Algorithm)
		if err != nil {
			fmt.Printf("Error encrypting message: %v\n", err)
			continue
		}
		
		// Decrypt message
		fmt.Println("Decrypting message...")
		decrypted, err := crypto.Decrypt(encrypted, keyPair.PrivateKey)
		if err != nil {
			fmt.Printf("Error decrypting message: %v\n", err)
			continue
		}
		
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
		decoys, err := crypto.GenerateCognitiveDecoyKeys(keyPair, decoyCount)
		if err != nil {
			fmt.Printf("Error generating decoy keys: %v\n", err)
			continue
		}
		
		fmt.Printf("Generated %d decoy keys\n", len(decoys))
		
		// Test decryption with decoy keys (should fail)
		fmt.Println("Attempting decryption with decoy keys (should fail):")
		for i, decoy := range decoys {
			// In a real implementation, this would fail because decoy keys can't decrypt
			decoyDecrypted, err := crypto.Decrypt(encrypted, decoy.PrivateKey)
			
			if err != nil {
				fmt.Printf("  Decoy #%d: Failed as expected with error: %v\n", i+1, err)
			} else {
				match := string(decoyDecrypted) == string(message)
				if match {
					fmt.Printf("  Decoy #%d: WARNING - Decryption succeeded (this would indicate a security issue in a real implementation)\n", i+1)
				} else {
					fmt.Printf("  Decoy #%d: Decryption produced incorrect data as expected\n", i+1)
				}
			}
		}
		
		fmt.Println()
	}
	
	fmt.Println("Demo completed!")
} 