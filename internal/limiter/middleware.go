package limiter

import (
	"context"
	"net/http"
	"strings"
)

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		ip := strings.Split(r.RemoteAddr, ":")[0]
		apiKey := r.Header.Get("API_KEY")

		if apiKey == "" && ip == "" {
			http.Error(w, "API_KEY or IP is required", http.StatusBadRequest)
			return
		}

		allowed, err := rl.Allow(ctx, apiKey, ip)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
