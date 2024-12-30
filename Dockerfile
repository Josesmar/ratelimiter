FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

EXPOSE 8080

CMD ["go", "run", "cmd/server/main.go"]
