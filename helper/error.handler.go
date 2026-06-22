package helper

import (
	"Microservice/data/response"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin" // Gunakan alias untuk menghindari konflik nama
	"go.uber.org/zap"
)

// == Unused Code
type CustomError struct {
	Code     int
	Message  string
	FileName string
	AtLine   int
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s (code: %d, file: %s, line: %d)", e.Message, e.Code, e.FileName, e.AtLine)
}

func ResponseError(ctx *gin.Context, err CustomError) {
	fmt.Printf("🔴 Error: %d - %s | %s at Line: %d\n", err.Code, err.Message, err.FileName, err.AtLine)
	errResponse := response.Response{
		Success: false,
		Code:    int(err.Code),
		Message: fmt.Sprintf("%s - %s", http.StatusText(err.Code), err.Message),
		Data:    nil,
	}

	ctx.JSON(
		GetErrorCode(err.Code), errResponse)
}

func GetFileAndLine(err error) (fileName string, atLine int) {
	if err == nil {
		return "", 0
	}

	// Get the program counter (PC) for the error
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "", 0
	}

	// Get the file and line information for the PC
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()

	return frame.File, frame.Line
}

// ----------------------------------------------------------------

func ErrorCatcher(err error, code int, message *string) *ErrorModel {
	var errorMessage string

	if message != nil {
		errorMessage = *message
	} else {
		errorMessage = http.StatusText(code)
	}

	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("🔴 - Error: %d - %s | File: %s (At Line: %d) | Raw Error: %s\n", code, errorMessage, file, line, err.Error())

	// Log Error in Zap (Production)
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	GetLogger().Error(
		"Error",
		zap.String("timestamp", timeStamp),
		zap.Int("code", code),
		zap.String("error_message", errorMessage),
		zap.String("file", file),
		zap.Int("line", line),
	)

	return &ErrorModel{
		Code:    code,
		Message: errorMessage,
	}
}

func ErrorLog(err error, code int, message *string) {
	var errorMessage string

	if message != nil {
		errorMessage = *message
	} else {
		errorMessage = http.StatusText(code)
	}

	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("🔴 - Error: %d - %s | File: %s (At Line: %d) | Raw Error: %s\n", code, errorMessage, file, line, err.Error())

	// Log Error in Zap (Production)
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	GetLogger().Error(
		"Error",
		zap.String("timestamp", timeStamp),
		zap.Int("code", code),
		zap.String("error_message", errorMessage),
		zap.String("file", file),
		zap.Int("line", line),
	)
}

type ErrorModel struct {
	Code    int
	Message string
}

// implementasikan interface error
func (e *ErrorModel) Error() string {
	return e.Message
}

func GetErrorCode(code int) int {
	switch code {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 402:
		return http.StatusPaymentRequired
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 405:
		return http.StatusMethodNotAllowed
	case 500:
		return http.StatusInternalServerError
	case 501:
		return http.StatusNotImplemented
	case 502:
		return http.StatusBadGateway
	case 503:
		return http.StatusServiceUnavailable
	case 504:
		return http.StatusGatewayTimeout
	case 505:
		return http.StatusHTTPVersionNotSupported
	default:
		return http.StatusInternalServerError
	}
}
