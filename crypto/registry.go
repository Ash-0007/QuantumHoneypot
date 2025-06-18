package crypto

import "fmt"

// Registry maintains a collection of crypto providers
type Registry struct {
	kemProviders       map[Algorithm]KEMProvider
	signatureProviders map[Algorithm]SignatureProvider
}

// NewRegistry creates a new crypto registry
func NewRegistry() *Registry {
	return &Registry{
		kemProviders:       make(map[Algorithm]KEMProvider),
		signatureProviders: make(map[Algorithm]SignatureProvider),
	}
}

// RegisterKEMProvider adds a KEM provider to the registry
func (r *Registry) RegisterKEMProvider(provider KEMProvider) {
	r.kemProviders[provider.Name()] = provider
}

// RegisterSignatureProvider adds a signature provider to the registry
func (r *Registry) RegisterSignatureProvider(provider SignatureProvider) {
	r.signatureProviders[provider.Name()] = provider
}

// GetKEMProvider retrieves a KEM provider by name
func (r *Registry) GetKEMProvider(alg Algorithm) (KEMProvider, error) {
	provider, exists := r.kemProviders[alg]
	if !exists {
		return nil, fmt.Errorf("KEM provider not found: %s", alg)
	}
	return provider, nil
}

// GetSignatureProvider retrieves a signature provider by name
func (r *Registry) GetSignatureProvider(alg Algorithm) (SignatureProvider, error) {
	provider, exists := r.signatureProviders[alg]
	if !exists {
		return nil, fmt.Errorf("signature provider not found: %s", alg)
	}
	return provider, nil
}

// DefaultRegistry creates a registry with all available providers
func DefaultRegistry() *Registry {
	registry := NewRegistry()
	
	// Register KEM providers
	registry.RegisterKEMProvider(NewMLKEM768Provider())
	registry.RegisterKEMProvider(NewECDHProvider())
	
	// Register signature providers
	registry.RegisterSignatureProvider(NewMLDSA65Provider())
	registry.RegisterSignatureProvider(NewECDSAProvider())
	
	return registry
} 