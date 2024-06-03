# Nome do executável
APP_NAME := goexpert-clean-architecture

# Diretório onde o Makefile está localizado
BASE_DIR := $(shell pwd)

# Docker compose command
DOCKER_COMPOSE := docker-compose

# Docker compose file
DOCKER_COMPOSE_FILE := docker-compose.yaml

# Docker command
DOCKER := docker

# Comandos relacionados ao Go
GO := go

# Configuração de cores para saída
COLOR_RESET := \033[0m
COLOR_ERROR := \033[0;31m
COLOR_SUCCESS := \033[0;32m
COLOR_INFO := \033[0;33m

CMD_DIR := ./cmd/ordersystem
PROTO_DIR := ./internal/infra/grpc/protofiles
PB_DIR := ./internal/infra/grpc/pb

# Comandos do Makefile
.PHONY: all build run fmt test cover docker-build docker-up docker-down docker-compose-up clean setup grpc

all: build

# Compila a aplicação localmente
build:
	@echo "$(COLOR_INFO)==> Compilando a aplicação$(COLOR_RESET)"
	$(GO) build -o $(APP_NAME) .

# Roda a aplicação localmente
run:
	@echo "$(COLOR_INFO)==> Rodando a aplicação localmente$(COLOR_RESET)"
	$(GO) run $(CMD_DIR)

# Formata o código
fmt:
	@echo "$(COLOR_INFO)==> Formatando o código$(COLOR_RESET)"
	$(GO) fmt ./...

# Roda os testes e gera o relatório de cobertura
test:
	@echo "$(COLOR_INFO)==> Rodando testes$(COLOR_RESET)"
	$(GO) generate -v ./...
	go-acc ./...

# Abre o relatório de cobertura
cover:
	@echo "$(COLOR_INFO)==> Abrindo relatório de cobertura$(COLOR_RESET)"
	$(GO) tool cover -html coverage.txt

# Executa o build usando Docker
docker-build:
	@echo "$(COLOR_INFO)==> Construindo a aplicação com Docker$(COLOR_RESET)"
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build

# Levanta os containers usando Docker Compose
docker-up:
	@echo "$(COLOR_INFO)==> Iniciando os containers com Docker Compose$(COLOR_RESET)"
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

# Derruba os containers usando Docker Compose
docker-down:
	@echo "$(COLOR_INFO)==> Derrubando os containers com Docker Compose$(COLOR_RESET)"
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

# Levanta os containers usando Docker Compose (alvo específico)
docker-compose-up:
	@echo "$(COLOR_INFO)==> Derrubando containers antigos e iniciando novos containers com Docker Compose$(COLOR_RESET)"
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

# Limpa arquivos temporários e cache
clean:
	@echo "$(COLOR_INFO)==> Limpando arquivos temporários$(COLOR_RESET)"
	$(GO) clean

# Alvos para ajudar no desenvolvimento local
setup:
	@echo "$(COLOR_INFO)==> Configurando o ambiente de desenvolvimento$(COLOR_RESET)"
	go get -u github.com/ory/go-acc
	go mod download
	go mod tidy

# Regra para gerar código gRPC
grpc:
	@echo "$(COLOR_INFO)==> Gerando código gRPC$(COLOR_RESET)"
	protoc --go_out=$(PB_DIR) --go-grpc_out=$(PB_DIR) -I=$(PROTO_DIR) $(PROTO_DIR)/*.proto

# Alvo padrão caso nenhum seja especificado
.DEFAULT_GOAL := all
