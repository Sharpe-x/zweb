package logger

import (
	"log"
	"time"
	"zweb"
)

// Logger 日志中间件
func Logger() zweb.HandlerFunc {
	return func(context *zweb.Context) {
		// start time
		t := time.Now()
		context.Next()
		log.Printf("[%d] %s in %v", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}
