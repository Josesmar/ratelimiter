# Rate Limiter em Go

## Visão Geral
Este projeto é um middleware de Rate Limiting desenvolvido em Go. Ele permite limitar o número de requisições permitidas a um serviço web com base em dois critérios:

1. **Token de Acesso**: Limita as requisições baseadas em um token de acesso informado no cabeçalho da requisição (`API_KEY`).

2. **IP Address**: Limita as requisiçoes baseadas em um ip address. No contexto do seu rate limiter, o IP de origem da requisição é automaticamente identificado pelo servidor através do cabeçalho RemoteAddr no HTTP. Ou seja, você não precisa explicitamente passar o IP na requisição — ele é derivado do cliente que realiza a conexão.

```bash
for i in {1..10}; do curl -X GET http://localhost:8080/; echo ""; done
```

O Rate Limiter utiliza o Redis para armazenar as informações de limite e bloqueio. Ele pode ser configurado através de variáveis de ambiente definidas em um arquivo `.env`.

---

## Configuração
### Dependências
Certifique-se de que os seguintes itens estejam instalados em sua máquina:
- Go (versão 1.19 ou superior)
- Docker e Docker Compose
- Redis

### Arquivo `.env`
Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
REDIS_ADDR=localhost:6379  # Endereço do servidor Redis
TOKEN_MAX_REQUESTS=10   # Limite máximo de requisições por token de acesso
IP_MAX_REQUESTS=5       # Limite máximo de requisição por ip de acesso
BAN_DURATION=5s         # Duração do bloqueio ao exceder o limite
```

---

## Executando o Projeto
### 1. Subir a aplicação com Docker Compose
Utilize o seguinte `docker-compose.yml` para subir o Redis:

```yaml
version: '3.9'
services:
  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    restart: always
```

Execute o comando:
```bash
docker-compose up -d
```

O servidor estará acessível na porta `8080`.

---

## Endpoints
### 1. `GET /`
Rota principal para testar o Rate Limiter. Caso o limite seja excedido, o servidor responderá com:

- **Status Code**: `429`
- **Mensagem**: `You have reached the maximum number of requests or actions allowed within a certain time frame`

### 2. Cabeçalho da requisição
Todas as requisições devem incluir o cabeçalho `API_KEY` ou um endereco de IP. Lembrando que o ID é identificado automaticamente.

```http
API_KEY: your-api-token
```

Caso o cabeçalho não seja enviado ou não seja identificado um endereço de ip, o servidor responderá com:

- **Status Code**: `400`
- **Mensagem**: `API_KEY or IP is required`

### 3. Manipulando o Cabeçalho X-Forwarded-For
Se o servidor está atrás de um proxy ou balanceador de carga, o cabeçalho X-Forwarded-For pode ser usado para identificar o IP real do cliente. Para simular isso:
```bash
curl -X GET http://localhost:8080/ -H "X-Forwarded-For: 192.168.1.1"
```
---

## Testando o Rate Limiter
### Requisições com `curl`
#### Teste de Limite Excedido
Envie multiplas requisições (mais do que o limite configurado no `.env`):

1. Requisição com IP
Simplesmente faça a requisição, pois o servidor já utiliza o IP do cliente automaticamente:
```bash
for i in {1..10}; do curl -X GET http://localhost:8080/; echo ""; done
```
- **Resposta esperada**: A partir da quinta requisição (baseado no exemplo), o servidor retornará `429 Too Many Requests`.

2. Requisição com Token (somente)
Para enviar somente o token no cabeçalho API_KEY, use:

```bash
for i in {1..10}; do curl -X GET http://localhost:8080/ -H "API_KEY: mytoken"; echo ""; done
```
- **Resposta esperada**: A partir da décima primeira requisição (baseado no exemplo), o servidor retornará `429 Too Many Requests`.

3. Requisição com IP e Token
Envie o token no cabeçalho API_KEY e simule um IP usando o cabeçalho X-Forwarded-For. Isso é útil se o seu servidor estiver configurado para lidar com proxies:
```bash
for i in {1..10}; do curl -X GET http://localhost:8080/ -H "API_KEY: mytoken" -H "X-Forwarded-For: 192.168.1.100"; echo ""; done
```
- **Resposta esperada**: Nesse caso o número máximo de request definida para o token ira sobrepor o número de requests definida para o IP.

## Estrutura do Projeto
```
ratelimiter/
├── internal/
│   └── limiter/        # Implementação do Rate Limiter
├── main.go             # Inicialização do servidor
├── .env                # Configurações de ambiente
└── go.mod              # Dependências do projeto
```

---