# Post-Quantum Cognitive Decoys - Usage Examples

This document provides practical examples of how to use the Post-Quantum Cognitive Decoys (PQCD) system.

## Command Line Examples

### Generate Key Pair

```bash
# Generate a Kyber key pair with 5 decoys
curl -X POST http://localhost:8082/api/keys/generate \
  -H "Content-Type: application/json" \
  -d '{"algorithm": "kyber", "count": 5}'
```

### Generate Decoys

```bash
# Generate 10 decoys for "kyber768" with complexity level 5
curl -X POST http://localhost:8082/api/decoys/generate \
  -H "Content-Type: application/json" \
  -d '{"target": "kyber768", "complexity": 5, "count": 10}'
```

### Encrypt a Message

```bash
# Encrypt a message using a public key
curl -X POST http://localhost:8082/api/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "plaintext": "This is a secret message",
    "public_key": "BASE64_ENCODED_PUBLIC_KEY",
    "algorithm": "kyber"
  }'
```

### Decrypt a Message

```bash
# Decrypt a message using a private key
curl -X POST http://localhost:8082/api/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "ciphertext": "BASE64_ENCODED_CIPHERTEXT",
    "private_key": "BASE64_ENCODED_PRIVATE_KEY",
    "nonce": "BASE64_ENCODED_NONCE",
    "algorithm": "kyber"
  }'
```

## JavaScript Examples

### Generate Key Pair

```javascript
async function generateKeyPair() {
  const response = await fetch('http://localhost:8082/api/keys/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      algorithm: 'kyber',
      count: 5
    }),
  });
  
  const data = await response.json();
  console.log('Generated key pair:', data);
  
  // Store the public key and fingerprint
  localStorage.setItem('publicKey', data.public_key);
  localStorage.setItem('keyFingerprint', data.fingerprint);
  
  return data;
}
```

### Generate Decoys

```javascript
async function generateDecoys(target, complexity = 5, count = 10) {
  const response = await fetch('http://localhost:8082/api/decoys/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      target,
      complexity,
      count
    }),
  });
  
  const data = await response.json();
  console.log('Generated decoys:', data);
  return data.decoys;
}
```

### Encrypt a Message

```javascript
async function encryptMessage(message, publicKey, algorithm = 'kyber') {
  const response = await fetch('http://localhost:8082/api/encrypt', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      plaintext: message,
      public_key: publicKey,
      algorithm
    }),
  });
  
  const data = await response.json();
  console.log('Encrypted message:', data);
  return {
    ciphertext: data.ciphertext,
    nonce: data.nonce,
    algorithm: data.algorithm
  };
}
```

### Decrypt a Message

```javascript
async function decryptMessage(ciphertext, privateKey, nonce, algorithm = 'kyber') {
  const response = await fetch('http://localhost:8082/api/decrypt', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      ciphertext,
      private_key: privateKey,
      nonce,
      algorithm
    }),
  });
  
  const data = await response.json();
  console.log('Decrypted message:', data);
  return data.plaintext;
}
```

## Python Examples

### Generate Key Pair

```python
import requests
import json
import base64

def generate_key_pair(algorithm="kyber", count=5):
    url = "http://localhost:8082/api/keys/generate"
    payload = {
        "algorithm": algorithm,
        "count": count
    }
    
    response = requests.post(url, json=payload)
    data = response.json()
    
    print(f"Generated {algorithm} key pair with fingerprint: {data['fingerprint']}")
    return data

# Example usage
key_data = generate_key_pair("kyber", 5)
public_key = key_data["public_key"]
```

### Generate Decoys

```python
def generate_decoys(target, complexity=5, count=10):
    url = "http://localhost:8082/api/decoys/generate"
    payload = {
        "target": target,
        "complexity": complexity,
        "count": count
    }
    
    response = requests.post(url, json=payload)
    data = response.json()
    
    print(f"Generated {len(data['decoys'])} decoys for target '{target}'")
    return data["decoys"]

# Example usage
decoys = generate_decoys("kyber768", 5, 10)
for i, decoy in enumerate(decoys):
    print(f"Decoy {i+1}: {decoy}")
```

### Encrypt a Message

```python
def encrypt_message(message, public_key, algorithm="kyber"):
    url = "http://localhost:8082/api/encrypt"
    payload = {
        "plaintext": message,
        "public_key": public_key,
        "algorithm": algorithm
    }
    
    response = requests.post(url, json=payload)
    data = response.json()
    
    print(f"Message encrypted using {algorithm}")
    return {
        "ciphertext": data["ciphertext"],
        "nonce": data["nonce"],
        "algorithm": data["algorithm"]
    }

# Example usage
message = "This is a secret message"
encrypted = encrypt_message(message, public_key)
```

### Decrypt a Message

```python
def decrypt_message(ciphertext, private_key, nonce, algorithm="kyber"):
    url = "http://localhost:8082/api/decrypt"
    payload = {
        "ciphertext": ciphertext,
        "private_key": private_key,
        "nonce": nonce,
        "algorithm": algorithm
    }
    
    response = requests.post(url, json=payload)
    data = response.json()
    
    print(f"Message decrypted using {algorithm}")
    return data["plaintext"]

# Example usage
decrypted_message = decrypt_message(
    encrypted["ciphertext"],
    private_key,
    encrypted["nonce"],
    encrypted["algorithm"]
)
print(f"Decrypted message: {decrypted_message}")
```

## Complete End-to-End Example

Here's a complete example showing the entire workflow:

```python
import requests
import json
import base64

# 1. Generate a key pair
def generate_key_pair():
    response = requests.post(
        "http://localhost:8082/api/keys/generate",
        json={"algorithm": "kyber", "count": 5}
    )
    return response.json()

# 2. Encrypt a message
def encrypt_message(message, public_key):
    response = requests.post(
        "http://localhost:8082/api/encrypt",
        json={
            "plaintext": message,
            "public_key": public_key,
            "algorithm": "kyber"
        }
    )
    return response.json()

# 3. Decrypt a message
def decrypt_message(ciphertext, private_key, nonce):
    response = requests.post(
        "http://localhost:8082/api/decrypt",
        json={
            "ciphertext": ciphertext,
            "private_key": private_key,
            "nonce": nonce,
            "algorithm": "kyber"
        }
    )
    return response.json()

# Execute the workflow
if __name__ == "__main__":
    # Step 1: Generate key pair
    print("Generating key pair...")
    key_data = generate_key_pair()
    public_key = key_data["public_key"]
    print(f"Key fingerprint: {key_data['fingerprint']}")
    
    # In a real scenario, the private key would be securely stored
    # For this example, we're simulating by retrieving it from the database
    
    # Step 2: Encrypt a message
    message = "This is a top secret message protected by post-quantum cryptography!"
    print(f"Encrypting message: {message}")
    encrypted = encrypt_message(message, public_key)
    
    print(f"Ciphertext: {encrypted['ciphertext'][:30]}...")
    print(f"Nonce: {encrypted['nonce']}")
    
    # Step 3: Decrypt the message
    # In a real scenario, you would retrieve the private key from secure storage
    # For this example, we're assuming we have access to it
    private_key = "RETRIEVED_PRIVATE_KEY"  # This would come from secure storage
    
    print("Decrypting message...")
    decrypted = decrypt_message(encrypted["ciphertext"], private_key, encrypted["nonce"])
    
    print(f"Decrypted message: {decrypted['plaintext']}")
    
    # Verify the decryption worked correctly
    if decrypted["plaintext"] == message:
        print("Success! The message was correctly encrypted and decrypted.")
    else:
        print("Error: The decrypted message doesn't match the original.") 