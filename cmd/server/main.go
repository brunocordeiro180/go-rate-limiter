package main

import (
	"log"
	"net/http"

	"github.com/brunocordeiro180/go-rate-limiter/config"
	"github.com/brunocordeiro180/go-rate-limiter/internal/infra/database"
	middleware "github.com/brunocordeiro180/go-rate-limiter/internal/infra/webserver"
	"github.com/brunocordeiro180/go-rate-limiter/internal/pkg/ratelimiter"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	storage := database.NewRedisConnection(cfg.RedisAddr, cfg.RedisPassword)
	limiter := ratelimiter.NewRateLimiter(storage, cfg.RateLimitIP, cfg.RateLimitToken, cfg.RateDuration)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	port := ":8080"

	handler := middleware.RateLimiterMiddleware(limiter, mux)

	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(port, handler))
}
