package s3presign

import (
	"strings"
	"testing"
	"time"
)

// TestPresignMatchesAWSReferenceVector verifies the SigV4 core against AWS's
// official "Example: signature calculation for presigned URL" from the S3 REST
// authentication docs. Matching this documented signature proves the canonical
// request, string-to-sign, signing-key derivation and encoding are all correct.
//
//	GET examplebucket/test.txt, region us-east-1, expires 86400s,
//	access key AKIAIOSFODNN7EXAMPLE, on 2013-05-24T00:00:00Z.
func TestPresignMatchesAWSReferenceVector(t *testing.T) {
	p := &Presigner{
		region:    "us-east-1",
		bucket:    "examplebucket",
		accessKey: "AKIAIOSFODNN7EXAMPLE",
		secretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}
	now := time.Date(2013, 5, 24, 0, 0, 0, 0, time.UTC)

	url, err := p.presign("GET", "test.txt", "examplebucket.s3.amazonaws.com", now, 86400*time.Second)
	if err != nil {
		t.Fatalf("presign returned error: %v", err)
	}

	const wantSig = "aeeed9bbccd4d02ee5c0109b86d86835f995330da4c265957d157751f604d404"
	if !strings.Contains(url, "X-Amz-Signature="+wantSig) {
		t.Fatalf("signature mismatch.\n got url: %s\nwant sig: %s", url, wantSig)
	}
}

func TestAWSURIEncode(t *testing.T) {
	cases := []struct {
		in          string
		encodeSlash bool
		want        string
	}{
		{"signatures/a b.png", false, "signatures/a%20b.png"},
		{"signatures/a b.png", true, "signatures%2Fa%20b.png"},
		{"keep-_.~", false, "keep-_.~"},
		{"a+b=c", true, "a%2Bb%3Dc"},
	}
	for _, c := range cases {
		if got := awsURIEncode(c.in, c.encodeSlash); got != c.want {
			t.Errorf("awsURIEncode(%q, %v) = %q, want %q", c.in, c.encodeSlash, got, c.want)
		}
	}
}

func TestNewReturnsNilWhenUnconfigured(t *testing.T) {
	if New("", "b", "k", "s") != nil {
		t.Error("expected nil Presigner when region is empty")
	}
	if New("r", "b", "k", "s") == nil {
		t.Error("expected non-nil Presigner when fully configured")
	}
	var p *Presigner
	if p.Enabled() {
		t.Error("nil Presigner should report Enabled() == false")
	}
}
