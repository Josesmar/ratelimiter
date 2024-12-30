package limiter

import (
	"context"
	"log"
	"ratelimiter/internal/limiter/persistence"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient      persistence.Store
	tokenMaxRequests int
	ipMaxRequest     int
	banDuration      time.Duration
}

func NewRateLimiter(redisClient persistence.Store, tokenMaxRequests, ipMaxRequest int, banDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient:      redisClient,
		tokenMaxRequests: tokenMaxRequests,
		ipMaxRequest:     ipMaxRequest,
		banDuration:      banDuration,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, apiKey, ip string) (bool, error) {
	var key string
	var maxRequests int

	if apiKey != "" {
		key = "token:" + apiKey
		maxRequests = rl.tokenMaxRequests
	} else {
		key = "ip:" + ip
		maxRequests = rl.ipMaxRequest
	}

	// Tenta obter o valor atual do contador
	count, err := rl.redisClient.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return false, err
	}

	// Se a chave não existir, inicialize com 1 e defina a expiração
	if count == 0 {
		err := rl.redisClient.Incr(ctx, key) // Incrementa para criar a chave
		if err != nil {
			return false, err
		}

		err = rl.redisClient.Expire(ctx, key, rl.banDuration) // Define a expiração
		if err != nil {
			return false, err
		}

		return true, nil
	}

	// Se o contador já atingiu o máximo permitido
	if count >= maxRequests {
		log.Printf("Limite excedido para %s: %d requisições", key, count)
		return false, nil
	}

	// Incrementa o contador, sem redefinir a expiração
	err = rl.redisClient.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	return true, nil
}
