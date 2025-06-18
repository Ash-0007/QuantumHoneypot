package main

import (
	"encoding/base64"
	"fmt"
	"log"
	
	"./crypto"
)

func main() {
	// Test Kyber key generation
	testAlgorithm("kyber")
}

func testAlgorithm(algorithm string) {
	fmt.Printf("\n=== Testing %s ===\n", algorithm)
	
	// Generate key pair
	fmt.Printf("Generating %s key pair...\n", algorithm)
	keyPair, err := crypto.GenerateKeyPair(algorithm)
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	
	// Print key information
	fmt.Printf("Public Key (%d bytes): %s\n", len(keyPair.PublicKey), base64.StdEncoding.EncodeToString(keyPair.PublicKey[:20])+"...")
	fmt.Printf("Private Key (%d bytes): %s\n", len(keyPair.PrivateKey), base64.StdEncoding.EncodeToString(keyPair.PrivateKey[:20])+"...")
	fmt.Printf("Fingerprint: %s\n", crypto.FingerPrint(keyPair.PublicKey))
	
	// Test message
	message := []byte("This is a test message for post-quantum encryption.")
	fmt.Printf("\nOriginal message: %s\n", string(message))
	
	// Encrypt
	fmt.Println("Encrypting message...")
	encrypted, err := crypto.Encrypt(message, keyPair.PublicKey, algorithm)
	if err != nil {
		log.Fatalf("Encryption failed: %v", err)
	}
	fmt.Printf("Ciphertext (%d bytes): %s\n", len(encrypted.Ciphertext), base64.StdEncoding.EncodeToString(encrypted.Ciphertext[:20])+"...")
	
	// Decrypt
	fmt.Println("Decrypting message...")
	decrypted, err := crypto.Decrypt(encrypted, keyPair.PrivateKey)
	if err != nil {
		log.Fatalf("Decryption failed: %v", err)
	}
	fmt.Printf("Decrypted message: %s\n", string(decrypted))
	
	// Verify
	if string(decrypted) == string(message) {
		fmt.Println("✅ Encryption and decryption successful!")
	} else {
		fmt.Println("❌ Decrypted message does not match original!")
	}
} 