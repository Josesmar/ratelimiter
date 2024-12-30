package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var secretKey = []byte("rate-limited-secret")

type TokenResponse struct {
	Token string `json:"token"`
}

func generateJWT() (string, error) {
	tokenMaxRequests, err := strconv.Atoi(os.Getenv("TOKEN_MAX_REQUESTS"))
	if err != nil {
		log.Printf("Erro ao converter TOKEN_MAX_REQUESTS: %v", err)
		return "", err
	}

	log.Printf("TOKEN_MAX_REQUESTS: %d", tokenMaxRequests)

	claims := jwt.MapClaims{
		"sub":  "1234567890",
		"name": "go-expert",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Duration(tokenMaxRequests) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Erro ao assinar o token: %v", err)
		return "", err
	}
	return signedToken, nil
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Iniciando a resposta para a rota padrão")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tudo Ok - Você pode continuar"))
}

func SetupRouter(r *mux.Router) {
	r.HandleFunc("/", welcomeHandler).Methods("GET")
}
