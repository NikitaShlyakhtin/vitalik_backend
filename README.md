# vitalik_backend

## About

This is a backend for CEX prototype, 
build for Innopolis University F24 FTCS course. 

The backend is developed with Go programming language and provides 
HTTP REST API for interaction with CEX: managing wallets, 
making deposits and trading orders.

## How to run

It's easy to run backend with PostgreSQL database with `docker-compose` CLI:

```shell
docker-compose up --build
```

Now, you are ready to make API requests. Use `localhost:8080` as a base URL.
