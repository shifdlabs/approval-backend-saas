package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// GetLogger initializes the logger only once
func GetLogger() *zap.Logger {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	})
	return logger
}

func PrintObject(obj interface{}, identifier string) {
	bytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println("Failed to print object:", err)
		return
	}

	fmt.Printf("🟠 %s: %s\n", identifier, string(bytes))
}

func PrintValue(value interface{}, identifier string) {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		fmt.Printf("🟠 %s: %v\n", identifier, value)
	default:
		fmt.Printf("🔵 %s (Unsupported Type): %v\n", identifier, value)
	}
}
