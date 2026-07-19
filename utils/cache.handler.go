package utils

import (
	"Microservice/config"
	"Microservice/helper"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

// scopedKey namespaces every cache entry by the caller's org_id so one tenant
// can never read another tenant's cached list (AUDIT: multi-tenant isolation).
// The empty org prefix is harmless because SubscriptionCheck rejects requests
// with no org_id before any cached route runs.
func scopedKey(ctx *gin.Context, key string) string {
	return helper.GetOrgID(ctx) + ":" + key
}

func GetCache(ctx *gin.Context, key string, dataType any) *any {
	cacheResult, _ := config.RedisClient.Get(ctx, scopedKey(ctx, key)).Result()
	response := dataType
	if err := json.Unmarshal([]byte(cacheResult), &response); err == nil {
		return &response
	} else {
		return nil
	}
}

func SetCache(ctx *gin.Context, key string, data any) {
	// Marshal the response
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Optionally handle error
		return
	}

	// Write under the caller-supplied key (previously hardcoded to
	// "All Documents", which collided across resource types), org-scoped.
	errRefresh := config.RedisClient.Set(ctx, scopedKey(ctx, key), jsonData, 60*time.Second).Err()
	if errRefresh != nil {
		println(errRefresh)
	}
}
