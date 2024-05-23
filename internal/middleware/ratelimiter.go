package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"rate-limiter-challenge-go/internal/storage"
	"time"
)

type RateLimiterRateConfig struct {
	MaxRequestsPerSecond  int64 `json:"maxRequestsPerSecond"`
	BlockTimeMilliseconds int64 `json:"blockTimeMilliseconds"`
}

type RateLimiterConfig struct {
	LimitByIP      *RateLimiterRateConfig
	LimitByToken   *RateLimiterRateConfig
	StorageAdapter storage.Adapter
	CustomTokens   *map[string]*RateLimiterRateConfig
}

func NewRateLimiter(config *RateLimiterConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return rateLimiter(config, next)
	}
}

func rateLimiter(config *RateLimiterConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		block, err := getBlock(config, r)

		if err != nil {
			handleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if block != nil {
			handleError(w, http.StatusTooManyRequests, "You have reached the maximum number of requests or actions allowed within a certain time frame.")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getBlock(config *RateLimiterConfig, r *http.Request) (*time.Time, error) {
	token := r.Header.Get("API_KEY")
	if token != "" {
		return handleToken(config, r, token)
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return checkRateLimit(r.Context(), "IP", host, config, config.LimitByIP)
}

func handleToken(config *RateLimiterConfig, r *http.Request, token string) (*time.Time, error) {
	var tokenConfig *RateLimiterRateConfig
	customTokenConfig, ok := (*config.CustomTokens)[token]
	if ok {
		tokenConfig = customTokenConfig
	} else {
		tokenConfig = config.LimitByToken
	}

	return checkRateLimit(r.Context(), "TOKEN", token, config, tokenConfig)
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}

func checkRateLimit(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
	if key == "" {
		return nil, nil
	}

	block, err := config.StorageAdapter.GetBlock(ctx, keyType, key)
	if err != nil {
		return nil, err
	}

	if block == nil {
		success, count, err := config.StorageAdapter.AddAccess(ctx, keyType, key, rateConfig.MaxRequestsPerSecond)
		if err != nil {
			return nil, err
		}

		if success {
			fmt.Printf("Access count within this window: %d\n", count)
		} else {
			fmt.Println("Access Denied")
			block, err = config.StorageAdapter.AddBlock(ctx, keyType, key, rateConfig.BlockTimeMilliseconds)
			if err != nil {
				return nil, err
			}
		}
	}

	if block != nil {
		fmt.Println("Blocked for", rateConfig.BlockTimeMilliseconds/1000, "seconds")
		return block, nil
	}

	return nil, nil
}
