package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ExtendWriteDeadline pushes back the connection's write deadline for a single
// route.
//
// The server sets WriteTimeout to 10s, which is right for every normal endpoint
// but far too short for AI draft generation: SAIS caps its own provider call at
// 60s, so a slow generation would have its connection torn down before the
// answer could be written. Rather than weakening WriteTimeout for the whole
// service, this raises the deadline only where it is needed.
//
// If the deadline cannot be set (e.g. a ResponseWriter that does not support
// it), the request proceeds under the server's default WriteTimeout.
func ExtendWriteDeadline(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := http.NewResponseController(c.Writer).SetWriteDeadline(time.Now().Add(d)); err != nil {
			log.Printf("middleware: could not extend write deadline for %s: %v", c.FullPath(), err)
		}
		c.Next()
	}
}
