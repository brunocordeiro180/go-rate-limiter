package middleware

import (
	"net/http"

	"github.com/brunocordeiro180/go-rate-limiter/internal/pkg/ratelimiter"
)

// RateLimiterMiddleware is a middleware that applies rate limiting logic to incoming requests.
// It checks if the request is over the limit based on the client IP or provided API_KEY.
// If the client has exceeded the limit, it returns HTTP 429 Too Many Requests.
// Otherwise, it forwards the request to the next handler.
func RateLimiterMiddleware(l *ratelimiter.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		isToken := false

		token := r.Header.Get("API_KEY")
		if token != "" {
			key = token
			isToken = true
		}

		if l.IsBlocked(r.Context(), key) || !l.Check(r.Context(), key, isToken) {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
