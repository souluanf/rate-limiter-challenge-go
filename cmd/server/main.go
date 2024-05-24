package main

import (
	"log"
	"net/http"
	"rate-limiter-challenge-go/config"
	"rate-limiter-challenge-go/internal/middleware"
	"rate-limiter-challenge-go/internal/storage"
	"rate-limiter-challenge-go/internal/storage/memory"
	"rate-limiter-challenge-go/internal/storage/redis"
	"rate-limiter-challenge-go/internal/webserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	webServer := webserver.NewWebServer(cfg.WebServerPort)
	var storageAdapter storage.Adapter
	if cfg.RedisAddr != "" {
		storageAdapter, err = redis.InitRedisAdapter(cfg.RedisAddr)
		if err != nil {
			log.Fatalf("Error initializing Redis adapter: %v", err)
		}
	} else {
		storageAdapter = memory.NewMemoryAdapter()
	}

	customTokens := populateCustomTokens()
	rateLimiterConfig := &middleware.RateLimiterConfig{
		LimitByIP: &middleware.RateLimiterRateConfig{
			MaxRequestsPerSecond:  cfg.LimitByIPMaxRPS,
			BlockTimeMilliseconds: cfg.LimitByIPBlockTimeMs,
		},
		LimitByToken: &middleware.RateLimiterRateConfig{
			MaxRequestsPerSecond:  cfg.LimitByTokenMaxRPS,
			BlockTimeMilliseconds: cfg.LimitByTokenBlockTimeMs,
		},
		StorageAdapter: storageAdapter,
		CustomTokens:   &customTokens,
	}

	rateLimiter := middleware.NewRateLimiter(rateLimiterConfig)
	webServer.Use(rateLimiter)
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			return
		}
	}
	webServer.AddHandler("/", rootHandler, "GET")
	webServer.Start()
}

func populateCustomTokens() map[string]*middleware.RateLimiterRateConfig {
	return map[string]*middleware.RateLimiterRateConfig{
		"ABC": {
			MaxRequestsPerSecond:  20,
			BlockTimeMilliseconds: 3000,
		},
		"DEF": {
			MaxRequestsPerSecond:  20,
			BlockTimeMilliseconds: 3000,
		},
	}
}
