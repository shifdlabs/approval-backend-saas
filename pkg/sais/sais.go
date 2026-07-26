// Package sais talks to the Shifd AI Service (SAIS), the shared Shifd Labs
// gateway for LLM features. SAIS holds the provider API key; the Approval
// Backend only forwards the end user's SIS token plus a shared internal secret
// that proves the caller is a trusted Shifd Labs backend.
package sais

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// productID identifies this product in SAIS's usage log, used for per-product
// cost accounting across Shifd Labs.
const productID = "shifd-approval"

// requestTimeout bounds one draft call. SAIS caps its own provider call at 60s
// and answers with an error after that, so this only needs a little headroom
// on top; it exists so a hung connection cannot pin an Approval request open.
const requestTimeout = 75 * time.Second

// maxErrorBodyBytes caps how much of a SAIS error body is read into an error
// message — enough to debug, not enough to dump a response into the logs.
const maxErrorBodyBytes = 512

// Client calls SAIS. Construct with NewClient; the zero value is not usable.
type Client struct {
	baseURL string
	secret  string
	http    *http.Client
}

func NewClient(baseURL, internalSecret string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		secret:  internalSecret,
		http:    &http.Client{Timeout: requestTimeout},
	}
}

// Enabled reports whether SAIS is configured. When it is not, AI endpoints
// answer with a clean "unavailable" instead of failing at call time — the rest
// of Approval must keep working without SAIS.
func (c *Client) Enabled() bool {
	return c != nil && c.baseURL != "" && c.secret != ""
}

// DraftRequest mirrors SAIS's POST /api/ai/draft body.
type DraftRequest struct {
	LetterType  string   `json:"letter_type"`
	Purpose     string   `json:"purpose"`
	KeyPoints   []string `json:"key_points"`
	CompanyName string   `json:"company_name,omitempty"`
	Recipient   string   `json:"recipient,omitempty"`
	Tone        string   `json:"tone,omitempty"`
}

type DraftResult struct {
	Draft      string `json:"draft"`
	TokensUsed int    `json:"tokens_used"`
}

// saisResponse is SAIS's envelope: {"data": {...}, "message": "success"}.
type saisResponse struct {
	Data DraftResult `json:"data"`
}

// GenerateDraft asks SAIS for a letter draft. bearerToken is the caller's raw
// Authorization header value — SAIS re-validates it against the same SIS JWKS,
// and derives the user/org it logs usage against from that token, so it must be
// forwarded unchanged rather than replaced by a service token.
func (c *Client) GenerateDraft(ctx context.Context, bearerToken string, req DraftRequest) (*DraftResult, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("sais: not configured")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("sais: failed to encode request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/ai/draft", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("sais: failed to build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", bearerToken)
	httpReq.Header.Set("X-Internal-Secret", c.secret)
	httpReq.Header.Set("X-Product-ID", productID)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sais: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		return nil, fmt.Errorf("sais: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(snippet)))
	}

	var parsed saisResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("sais: failed to decode response: %w", err)
	}
	if strings.TrimSpace(parsed.Data.Draft) == "" {
		return nil, fmt.Errorf("sais: response contained an empty draft")
	}

	return &parsed.Data, nil
}
