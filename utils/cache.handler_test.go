package utils

import (
	"Microservice/helper"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func ctxWithOrg(org string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(helper.ContextKeyOrgID, org)
	return c
}

// Guards the multi-tenant cache-isolation fix: the same logical key must map to
// distinct Redis keys per org (no cross-tenant read), and distinct logical keys
// must not collide within one org (the old code hardcoded "All Documents").
func TestScopedKeyIsolation(t *testing.T) {
	orgA := scopedKey(ctxWithOrg("org-a"), "All Documents")
	orgB := scopedKey(ctxWithOrg("org-b"), "All Documents")
	if orgA == orgB {
		t.Fatalf("cross-tenant collision: org-a and org-b share cache key %q", orgA)
	}

	docs := scopedKey(ctxWithOrg("org-a"), "All Documents")
	attach := scopedKey(ctxWithOrg("org-a"), "All Attachment")
	if docs == attach {
		t.Fatalf("resource-type collision within one org: %q", docs)
	}
}
