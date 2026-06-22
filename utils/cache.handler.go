package utils

import (
	"Microservice/config"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCache(ctx *gin.Context, key string, dataType any) *any {
	cacheResult, _ := config.RedisClient.Get(ctx, key).Result()
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

	errRefresh := config.RedisClient.Set(ctx, "All Documents", jsonData, 60*time.Second).Err()
	if errRefresh != nil {
		println(errRefresh)
	}
}
