package middleware

import (
	"jayant-001/api-gateway/internal/config"
	"net/http"
	"strconv"

	"golang.org/x/time/rate"
)

func RateLimit(cfg config.RateLimitConfig) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(cfg.RequestsPerSecond), cfg.Burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				// Add headers to the response when the request is denied
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(cfg.Burst))
				w.Header().Set("X-RateLimit-Remaining", "0")
				
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			// Add headers to the response for allowed requests
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(cfg.Burst))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limiter.Burst()-1)) // Tokens remaining after this request
			
			next.ServeHTTP(w, r)
		})
	}
}
