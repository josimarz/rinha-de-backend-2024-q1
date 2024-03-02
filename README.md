# Rinha de Backend - 2024/Q1

Este projeto foi criado para o desafio [**Rinha de Backend - 2024/Q1**](https://github.com/zanfranceschi/rinha-de-backend-2024-q1).

## Tecnologias usadas

### Linguagem de programação

Para desenvolvimento do código da aplicação foi utilizada a linguagem de programação Go. Go ou Golang, é uma linguagem fortemente tipada e compilada, o que permite a criação de programas com excelente desempenho e baixo consumo de memória.

### Banco de dados

Para armazenamento dos dados foi utilizado o Postgres, um SGBD (Sistema Gerenciador de Banco de Dados) relacional muito popular.

### Servidor HTTP

Para balanceamento de carga das requisições foi utilizado o Nginx.

### Controle de concorrência

Para realizar o controle de concorrência foi utilizado o banco de dados em memória Redis, com auxílio do algorítmo [Redlock](https://redis.io/docs/manual/patterns/distributed-locks/), implementado pela biblioteca [redsync](https://github.com/go-redsync/redsync).

## Como executar

Para executar o projeto, acesse o diretório raiz deste projeto e execute:

```sh
docker compose up -d
```