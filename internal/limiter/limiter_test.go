package limiter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"ratelimiter/internal/limiter/persistence"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func createTestLimiter(redisAddr string, maxRequests int, banDuration time.Duration) *RateLimiter {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := redisClient.FlushAll(context.Background()).Err(); err != nil {
		panic("Falha ao limpar Redis: " + err.Error())
	}

	redisStore := persistence.NewRedisStore(redisClient)
	return NewRateLimiter(redisStore, maxRequests, maxRequests, banDuration)
}

func TestRateLimiterByIP(t *testing.T) {
	limiter := createTestLimiter("localhost:6379", 5, 1*time.Minute)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := httptest.NewServer(limiter.Middleware(mux))
	defer server.Close()

	client := &http.Client{}
	ip := "192.168.1.1"
	for i := 1; i <= 7; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = ip

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro ao fazer a requisição: %v", err)
		}

		if i <= 5 && resp.StatusCode != http.StatusOK {
			t.Errorf("Esperava status 200 OK para requisição %d, mas recebeu %d", i, resp.StatusCode)
		} else if i > 5 && resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Esperava status 429 Too Many Requests para requisição %d, mas recebeu %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestRateLimiterByToken(t *testing.T) {
	limiter := createTestLimiter("localhost:6379", 10, 1*time.Minute)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := httptest.NewServer(limiter.Middleware(mux))
	defer server.Close()

	client := &http.Client{}
	token := "mytoken"
	for i := 1; i <= 12; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Add("API_KEY", token)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro ao fazer a requisição: %v", err)
		}

		if i <= 10 && resp.StatusCode != http.StatusOK {
			t.Errorf("Esperava status 200 OK para requisição %d, mas recebeu %d", i, resp.StatusCode)
		} else if i > 10 && resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Esperava status 429 Too Many Requests para requisição %d, mas recebeu %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestBanDuration(t *testing.T) {
	limiter := createTestLimiter("localhost:6379", 5, 5*time.Second)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := httptest.NewServer(limiter.Middleware(mux))
	defer server.Close()

	client := &http.Client{}
	ip := "192.168.1.1"
	for i := 1; i <= 6; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = ip

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro ao fazer a requisição: %v", err)
		}

		if i <= 5 && resp.StatusCode != http.StatusOK {
			t.Errorf("Esperava status 200 OK para requisição %d, mas recebeu %d", i, resp.StatusCode)
		} else if i == 6 && resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Esperava status 429 Too Many Requests para requisição %d, mas recebeu %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}

	time.Sleep(6 * time.Second)

	req, _ := http.NewRequest("GET", server.URL, nil)
	req.RemoteAddr = ip
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro ao fazer a requisição: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Esperava status 200 OK após banimento, mas recebeu %d", resp.StatusCode)
	}
	resp.Body.Close()
}
