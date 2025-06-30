package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

func ErrorJSONGin(c *gin.Context, status int, msg string) {
	response := ErrorResponse{
		Error: msg,
		Code:  fmt.Sprintf("ERR_%d", status),
	}
	c.JSON(status, response)
}

func ErrorJSONGinWithDetails(c *gin.Context, status int, msg, details string) {
	response := ErrorResponse{
		Error:   msg,
		Code:    fmt.Sprintf("ERR_%d", status),
		Details: details,
	}
	c.JSON(status, response)
}

func GinPrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		path := c.FullPath()

		if path == "" {
			path = c.Request.URL.Path
		}

		HTTPRequestDuration.With(prometheus.Labels{
			"path":        path,
			"method":      c.Request.Method,
			"status_code": fmt.Sprintf("%d", c.Writer.Status()),
		}).Observe(duration.Seconds())
	}
}

func MethodGuardGin(allowedMethod string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != allowedMethod {
			ErrorJSONGin(c, 405, ErrMethodNotAllowed)
			c.Abort()
		}
	}
}

// RateLimitMiddleware is a placeholder for rate limiting implementation
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement rate limiting logic
		// This could use a token bucket or sliding window approach
		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
