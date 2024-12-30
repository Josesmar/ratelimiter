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
	log.Println("Tentando carregar o arquivo .env")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env:", err)
	} else {
		log.Println("Arquivo .env carregado com sucesso!")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR não configurado no arquivo .env")
	}

	tokenMaxRequests, err := strconv.Atoi(os.Getenv("TOKEN_MAX_REQUESTS"))
	if err != nil {
		log.Fatal("Erro ao ler TOKEN_MAX_REQUESTS no arquivo .env:", err)
	}

	ipMaxRequest, err := strconv.Atoi(os.Getenv("IP_MAX_REQUESTS"))
	if err != nil {
		log.Fatal("Erro ao ler IP_MAX_REQUESTS no arquivo.env:", err)
	}

	banDuration, err := time.ParseDuration(os.Getenv("BAN_DURATION"))
	if err != nil {
		log.Fatal("Erro ao ler BAN_DURATION no arquivo .env:", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar com o Redis: %v", err)
	}
	log.Printf("Conectado ao Redis com sucesso: %s", pong)

	redisStore := persistence.NewRedisStore(redisClient)
	rateLimiter := limiter.NewRateLimiter(redisStore, tokenMaxRequests, ipMaxRequest, banDuration)

	r := mux.NewRouter()
	server.SetupRouter(r)
	handler := rateLimiter.Middleware(r)

	log.Println("Servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
