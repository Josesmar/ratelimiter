package limiter

// import (
// 	"context"
// 	"time"

// 	"github.com/go-redis/redis/v8"
// )

// type Storage interface {
// 	Increment(ctx context.Context, key string) (int, error)
// 	SetExpiration(ctx context.Context, key string, duration time.Duration) error
// }

// type RedisStorage struct {
// 	Client *redis.Client
// }

// func NewRedisStorage(addr string) (*RedisStorage, error) {
// 	client := redis.NewClient(&redis.Options{Addr: addr})
// 	if err := client.Ping(context.Background()).Err(); err != nil {
// 		return nil, err
// 	}
// 	return &RedisStorage{Client: client}, nil
// }

// func (r *RedisStorage) Increment(ctx context.Context, key string) (int, error) {
// 	val, err := r.Client.Incr(ctx, key).Result()
// 	return int(val), err
// }

// func (r *RedisStorage) SetExpiration(ctx context.Context, key string, duration time.Duration) error {
// 	return r.Client.Expire(ctx, key, duration).Err()
// }
