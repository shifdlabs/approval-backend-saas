// Package jwks fetches and caches the RSA public keys published by the Shifd
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

// JWKSClient holds the SIS JWKS URL, expected issuer, and the currently-cached
// RSA public keys indexed by `kid`. The zero value is not usable — construct
// it with NewJWKSClient.
type JWKSClient struct {
	url            string
	expectedIssuer string
	// keys maps kid → RSA public key. The empty-string key "" holds the first
	// encountered key as a fallback for tokens that omit `kid`.
	keys map[string]*rsa.PublicKey
	mu   sync.RWMutex
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

// NewJWKSClient creates a JWKS client and attempts an initial fetch of the SIS
// public keys. expectedIssuer is validated on every token (pass "" to skip).
// If SIS is unreachable (e.g. during local development), the client starts
// without cached keys and logs a warning — StartAutoRefresh keeps retrying.
func NewJWKSClient(jwksURL string, expectedIssuer string) (*JWKSClient, error) {
	c := &JWKSClient{url: jwksURL, expectedIssuer: expectedIssuer}
	if err := c.Refresh(); err != nil {
		log.Printf("jwks: could not fetch/parse JWKS from %q on startup (will retry in background): %v", jwksURL, err)
	}
	return c, nil
}

// GetPublicKey returns the cached RSA public key for the given kid. If kid is
// not found, it returns the fallback (first-encountered) key. Thread-safe.
func (j *JWKSClient) GetPublicKey(kid string) *rsa.PublicKey {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if key, ok := j.keys[kid]; ok {
		return key
	}
	return j.keys[""] // fallback for tokens without kid or unknown kid
}

// GetExpectedIssuer returns the issuer string this client validates against.
func (j *JWKSClient) GetExpectedIssuer() string {
	return j.expectedIssuer
}

// HasAnyKey reports whether the client has at least one cached key.
func (j *JWKSClient) HasAnyKey() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return len(j.keys) > 0
}

// Refresh re-fetches the JWKS and atomically swaps in the new key set.
// On any error the previously-cached keys are left untouched.
func (j *JWKSClient) Refresh() error {
	keys, err := fetchAndParse(j.url)
	if err != nil {
		return err
	}
	j.mu.Lock()
	j.keys = keys
	j.mu.Unlock()
	return nil
}

// StartAutoRefresh launches a background goroutine that calls Refresh on the
// given interval (handles SIS key rotation). A non-positive interval defaults
// to 24h. Refresh failures are logged and the last good keys keep serving.
// If no keys are cached yet (SIS was unreachable at startup), it retries every
// 10 s until the keys are obtained before switching to the normal interval.
func (j *JWKSClient) StartAutoRefresh(interval time.Duration) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	go func() {
		// Fast-retry loop: keep trying every 10 s until we have a key.
		if !j.HasAnyKey() {
			retry := time.NewTicker(10 * time.Second)
			for range retry.C {
				if err := j.Refresh(); err != nil {
					log.Printf("jwks: startup retry failed (will retry in 10s): %v", err)
				} else {
					log.Printf("jwks: public keys loaded from %s", j.url)
					retry.Stop()
					break
				}
			}
		}
		// Normal rotation interval.
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := j.Refresh(); err != nil {
				log.Printf("jwks: background refresh failed (keeping cached keys): %v", err)
			}
		}
	}()
}

// fetchAndParse downloads the JWKS document and returns all usable RSA signing
// keys indexed by kid. The empty-string key "" holds the first key as a
// fallback for tokens that omit the kid header.
func fetchAndParse(url string) (map[string]*rsa.PublicKey, error) {
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

	keys := make(map[string]*rsa.PublicKey)
	for _, k := range set.Keys {
		if !strings.EqualFold(k.Kty, "RSA") {
			continue
		}
		// "use" is optional; if present it must be "sig".
		if k.Use != "" && !strings.EqualFold(k.Use, "sig") {
			continue
		}
		pub, err := parseRSAPublicKey(k.N, k.E)
		if err != nil {
			log.Printf("jwks: skipping key %q: %v", k.Kid, err)
			continue
		}
		keys[k.Kid] = pub
		// Store the first valid key under "" as a fallback for kid-less tokens.
		if _, hasDefault := keys[""]; !hasDefault {
			keys[""] = pub
		}
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no usable RSA signing key found in JWKS (%d keys total)", len(set.Keys))
	}
	return keys, nil
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
