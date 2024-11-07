package middleware

import (
	"fmt"
	"log"
	"metrika/pkg/logger"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler, logDirect logger.LoggerType) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		duration := time.Since(start)

		ip := r.RemoteAddr
		userAgent := r.Header.Get("User-Agent")

		logMessage := fmt.Sprintf(
			"Request: %s %s from %s (User-Agent: %s) took %s",
			r.Method, r.URL.String(), ip, userAgent, duration.String(),
		)

		if err := logDirect.LoggerBasic(logger.INFO_LOG, logMessage); err != nil {
			log.Println("Failed to log to database:", err)
		}
	})

}
