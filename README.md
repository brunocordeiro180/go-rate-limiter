
# go-rate-limiter

Este projeto é uma implementação de um *rate limiter* em Go, desenvolvido como parte de um desafio técnico. Ele permite limitar o número de requisições por segundo com base no endereço IP ou em um token de acesso fornecido no cabeçalho da requisição.

## 🚀 Funcionalidades

- Limitação de requisições por IP ou por token de acesso.
- Configuração de limites e tempos de bloqueio via variáveis de ambiente ou arquivo `.env`.
- Middleware para fácil integração com servidores HTTP.
- Persistência dos dados de limitação utilizando Redis.
- Estrutura modular que permite substituição do mecanismo de persistência.
- Resposta com código HTTP 429 quando o limite é excedido.

## 📦 Estrutura do Projeto

```
go-rate-limiter/
├── cmd/
│   └── server/         
├── config/  
├── docs/             
├── internal/
│   ├── infra/
|       └── database/
|       └── webserver/
│   ├── pkg/      
├── docker-compose.yaml
├── .env           
├── header.lua
├── go.mod
└── go.sum
```

## ⚙️ Configuração

As configurações podem ser definidas via variáveis de ambiente ou no arquivo `.env` na raiz do projeto.

Variáveis disponíveis:

- `RATE_LIMIT_IP`: Número máximo de requisições por segundo por IP (ex: `10`).
- `RATE_LIMIT_TOKEN`: Número máximo de requisições por segundo por token (ex: `100`).
- `RATE_DURATION`: Tempo de bloqueio em segundos após exceder o limite (ex: `300` para 5 minutos).
- `REDIS_HOST`: Endereço do Redis (ex: `localhost:6379`).
- `REDIS_PASSWORD`: Senha do Redis, se aplicável.

## 🧪 Testes de Estresse com wrk

Para realizar testes de estresse, você pode utilizar a ferramenta `wrk`.

### Instalação do wrk

No macOS:

```bash
brew install wrk
```

No Linux:

```bash
sudo apt-get install wrk
```

### Execução Básica

Execute o seguinte comando para iniciar um teste de estresse:

```bash
wrk -t12 -c400 -d1s http://localhost:8080/
```

- `-t12`: 12 threads
- `-c400`: 400 conexões simultâneas
- `-d1s`: duração de 1 segundo

### Uso do header.lua

Para adicionar cabeçalhos personalizados, como o `API_KEY`, utilize o script `header.lua` incluído no projeto:

```bash
wrk -t12 -c400 -d1s -s header.lua http://localhost:8080/
```

O script `header.lua` adiciona automaticamente o cabeçalho `API_KEY` com um token de acesso.

## 🐳 Utilizando Docker Compose

Para iniciar o Redis utilizando Docker Compose:

```bash
docker-compose up -d
```

Certifique-se de que o Redis esteja em execução antes de iniciar o servidor Go.

## 🛠️ Executando o Servidor

Com o Redis em execução e as variáveis de ambiente configuradas, inicie o servidor:

```bash
go run cmd/server/main.go
```

O servidor estará disponível em `http://localhost:8080/`.