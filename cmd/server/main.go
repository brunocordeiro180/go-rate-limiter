package main

import (
	"log"
	"net/http"

	"github.com/brunocordeiro180/go-rate-limiter/config"
	_ "github.com/brunocordeiro180/go-rate-limiter/docs"
	"github.com/brunocordeiro180/go-rate-limiter/internal/infra/database"
	middleware "github.com/brunocordeiro180/go-rate-limiter/internal/infra/webserver"
	"github.com/brunocordeiro180/go-rate-limiter/internal/pkg/ratelimiter"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Rate Limiter API
// @version 1.0
// @description A basic rate-limiting API using Redis and native net/http.
// @termsOfService http://swagger.io/terms/

// @contact.name Bruno Cordeiro
// @contact.url https://github.com/brunocordeiro180
// @contact.email brunocordeiro180@gmail.com

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	storage := database.NewRedisConnection(cfg.RedisAddr, cfg.RedisPassword)
	limiter := ratelimiter.NewRateLimiter(storage, cfg.RateLimitIP, cfg.RateLimitToken, cfg.RateDuration)

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	port := ":8080"

	handler := middleware.RateLimiterMiddleware(limiter, mux)

	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(port, handler))
}

// @Summary Hello World Endpoint
// @Description Simple hello world endpoint with rate limiting applied
// @Tags hello-world
// @Produce plain
// @Success 200 {string} string "Hello, World!"
// @Failure 429 {string} string "Too Many Requests"
// @Router / [get]
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
