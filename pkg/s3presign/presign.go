// Package s3presign generates AWS Signature Version 4 presigned URLs for S3
// object operations (PUT/DELETE) without depending on the AWS SDK.
//
// Rationale (AUDIT SEC-01): before Phase 2 the frontend embedded long-lived IAM
// credentials in its bundle to upload directly to S3, leaking them to every
// browser. Now the Approval Backend holds the credentials and hands the browser
// a short-lived presigned URL; the browser PUTs/DELETEs straight to S3 with no
// credentials of its own.
package s3presign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	algorithm       = "AWS4-HMAC-SHA256"
	service         = "s3"
	aws4Request     = "aws4_request"
	unsignedPayload = "UNSIGNED-PAYLOAD"
)

// Presigner signs S3 URLs for a single bucket/region using static credentials.
type Presigner struct {
	region    string
	bucket    string
	accessKey string
	secretKey string
}

// New builds a Presigner. It returns nil when any required field is empty, so a
// nil *Presigner is a valid "S3 uploads are not configured" sentinel that
// callers can check with Enabled().
func New(region, bucket, accessKey, secretKey string) *Presigner {
	if region == "" || bucket == "" || accessKey == "" || secretKey == "" {
		return nil
	}
	return &Presigner{region: region, bucket: bucket, accessKey: accessKey, secretKey: secretKey}
}

// Enabled reports whether the presigner is configured (non-nil with creds).
func (p *Presigner) Enabled() bool { return p != nil }

// Presign returns a presigned URL for the given HTTP method ("PUT"/"DELETE") and
// object key, valid for the given duration. The payload is signed as
// UNSIGNED-PAYLOAD and only the host header is signed, so the browser is free to
// set Content-Type (and any other header) without invalidating the signature.
func (p *Presigner) Presign(method, key string, expires time.Duration) (string, error) {
	if !p.Enabled() {
		return "", fmt.Errorf("s3presign: not configured")
	}
	host := fmt.Sprintf("%s.s3.%s.amazonaws.com", p.bucket, p.region)
	return p.presign(strings.ToUpper(method), key, host, time.Now().UTC(), expires)
}

// presign is the deterministic core of Presign, split out so tests can pin the
// host and timestamp and check the signature against the AWS reference vector.
func (p *Presigner) presign(method, key, host string, now time.Time, expires time.Duration) (string, error) {
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	credentialScope := strings.Join([]string{dateStamp, p.region, service, aws4Request}, "/")

	// Canonical URI: single URI-encoding of the path, slashes preserved (S3 does
	// not double-encode). Keys never carry a leading slash, so add one here.
	canonicalURI := "/" + awsURIEncode(key, false)

	// Canonical query string: every SigV4 auth parameter, URI-encoded and sorted.
	params := map[string]string{
		"X-Amz-Algorithm":     algorithm,
		"X-Amz-Credential":    p.accessKey + "/" + credentialScope,
		"X-Amz-Date":          amzDate,
		"X-Amz-Expires":       strconv.Itoa(int(expires.Seconds())),
		"X-Amz-SignedHeaders": "host",
	}
	canonicalQuery := encodeCanonicalQuery(params)

	canonicalHeaders := "host:" + host + "\n"
	signedHeaders := "host"

	canonicalRequest := strings.Join([]string{
		method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		signedHeaders,
		unsignedPayload,
	}, "\n")

	stringToSign := strings.Join([]string{
		algorithm,
		amzDate,
		credentialScope,
		hexSHA256([]byte(canonicalRequest)),
	}, "\n")

	signingKey := deriveSigningKey(p.secretKey, dateStamp, p.region)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	return fmt.Sprintf("https://%s%s?%s&X-Amz-Signature=%s", host, canonicalURI, canonicalQuery, signature), nil
}

// encodeCanonicalQuery URI-encodes and sorts query parameters per SigV4. Keys
// are sorted by their encoded form; all keys here are pure ASCII, so encoded
// order matches raw order.
func encodeCanonicalQuery(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, awsURIEncode(k, true))
	}
	sort.Strings(keys)

	// Re-map encoded key → encoded value for assembly.
	encoded := make(map[string]string, len(params))
	for k, v := range params {
		encoded[awsURIEncode(k, true)] = awsURIEncode(v, true)
	}

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+encoded[k])
	}
	return strings.Join(parts, "&")
}

// awsURIEncode percent-encodes per RFC 3986 as AWS expects: the unreserved set
// (A-Z a-z 0-9 - _ . ~) is left as-is, everything else is uppercase %-encoded.
// Slash is encoded only when encodeSlash is true (query values yes, path no).
func awsURIEncode(s string, encodeSlash bool) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~':
			b.WriteByte(c)
		case c == '/' && !encodeSlash:
			b.WriteByte('/')
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func deriveSigningKey(secret, dateStamp, region string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	return hmacSHA256(kService, aws4Request)
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func hexSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
