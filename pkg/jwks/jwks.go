// Package jwks fetches and caches the RSA public key published by the Shifd
// Labs Identity Service (SIS) JWKS endpoint, so the Approval Backend can verify
// SIS-issued JWTs locally (offline) without calling SIS on every request.
package jwks

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

// JWKSClient holds the SIS JWKS URL and the currently-cached RSA public key.
// The zero value is not usable — construct it with NewJWKSClient.
type JWKSClient struct {
	url       string
	publicKey *rsa.PublicKey
	mu        sync.RWMutex
}

// jwk is a single JSON Web Key as published at the JWKS endpoint.
type jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"` // modulus    (base64url)
	E   string `json:"e"` // public exp (base64url)
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

// NewJWKSClient fetches the JWKS from jwksURL, parses the RSA public key, and
// caches it. The fetch happens synchronously on startup; if it fails we panic
// with a clear message because the service cannot authenticate anyone without
// SIS's public key.
func NewJWKSClient(jwksURL string) (*JWKSClient, error) {
	c := &JWKSClient{url: jwksURL}
	if err := c.Refresh(); err != nil {
		panic(fmt.Sprintf("jwks: could not fetch/parse JWKS from %q on startup: %v", jwksURL, err))
	}
	return c, nil
}

// GetPublicKey returns the cached RSA public key. Thread-safe.
func (j *JWKSClient) GetPublicKey() *rsa.PublicKey {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.publicKey
}

// Refresh re-fetches the JWKS and atomically swaps in the new public key.
// On any error the previously-cached key is left untouched.
func (j *JWKSClient) Refresh() error {
	key, err := fetchAndParse(j.url)
	if err != nil {
		return err
	}
	j.mu.Lock()
	j.publicKey = key
	j.mu.Unlock()
	return nil
}

// StartAutoRefresh launches a background goroutine that calls Refresh on the
// given interval (handles SIS key rotation). A non-positive interval defaults
// to 24h. Refresh failures are logged and the last good key keeps serving.
func (j *JWKSClient) StartAutoRefresh(interval time.Duration) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := j.Refresh(); err != nil {
				log.Printf("jwks: background refresh failed (keeping cached key): %v", err)
			}
		}
	}()
}

// fetchAndParse downloads the JWKS document and returns the first usable RSA
// signing key.
func fetchAndParse(url string) (*rsa.PublicKey, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch jwks: unexpected status %d", resp.StatusCode)
	}

	var set jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, fmt.Errorf("decode jwks: %w", err)
	}

	for _, k := range set.Keys {
		if !strings.EqualFold(k.Kty, "RSA") {
			continue
		}
		// "use" is optional; if present it must be "sig".
		if k.Use != "" && !strings.EqualFold(k.Use, "sig") {
			continue
		}
		return parseRSAPublicKey(k.N, k.E)
	}
	return nil, fmt.Errorf("no usable RSA signing key found in JWKS (%d keys)", len(set.Keys))
}

// parseRSAPublicKey builds an *rsa.PublicKey from the base64url modulus (n) and
// exponent (e) fields of a JWK.
func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := decodeBase64URL(nStr)
	if err != nil {
		return nil, fmt.Errorf("decode modulus: %w", err)
	}
	eBytes, err := decodeBase64URL(eStr)
	if err != nil {
		return nil, fmt.Errorf("decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)

	var e int
	for _, b := range eBytes {
		e = e<<8 | int(b)
	}
	if e == 0 {
		e = 65537 // sensible default if exponent was omitted/zero
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}

// decodeBase64URL decodes a base64url string, tolerating both padded and
// unpadded variants (JWKS values are conventionally unpadded).
func decodeBase64URL(s string) ([]byte, error) {
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(s)
}
