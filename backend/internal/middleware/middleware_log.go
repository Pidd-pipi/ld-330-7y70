package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

func RequestLogger(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Header("X-Request-ID", rid)
		start := time.Now()
		c.Next()
		l.Info("request", "request_id", rid, "method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status(), "latency_ms", time.Since(start).Milliseconds())
	}
}
