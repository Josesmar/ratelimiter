package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"ratelimiter/internal/limiter"
	"ratelimiter/internal/limiter/persistence"
	"ratelimiter/internal/server"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("Trying to load the .env file")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	} else {
		log.Println(".env file loaded successfully!")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR not configured in the .env file")
	}

	tokenMaxRequests, err := strconv.Atoi(os.Getenv("TOKEN_MAX_REQUESTS"))
	if err != nil {
		log.Fatal("Error reading TOKEN_MAX_REQUESTS from .env file:", err)
	}

	ipMaxRequest, err := strconv.Atoi(os.Getenv("IP_MAX_REQUESTS"))
	if err != nil {
		log.Fatal("Error reading IP_MAX_REQUESTS from .env file:", err)
	}

	banDuration, err := time.ParseDuration(os.Getenv("BAN_DURATION"))
	if err != nil {
		log.Fatal("Error reading BAN_DURATION from .env file:", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Printf("Successfully connected to Redis: %s", pong)

	redisStore := persistence.NewRedisStore(redisClient)
	rateLimiter := limiter.NewRateLimiter(redisStore, tokenMaxRequests, ipMaxRequest, banDuration)

	r := mux.NewRouter()
	server.SetupRouter(r)
	handler := rateLimiter.Middleware(r)

	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
