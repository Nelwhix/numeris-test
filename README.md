# Numeris-Test Api
Stack:
- Go
- Postgres
- Redis

## Setup
start up the Postgres and Redis containers by running 
```bash
    docker compose up -d
```

Copy the .env with:
```bash
 cp .env.example .env
```
All the necessary creds are in the compose.yml, you can edit as you wish

Download dependencies with:
```bash
    go mod download
```

Run migrations with:
```bash 
 tern migrate --migrations ./migrations
```
To download tern: https://github.com/jackc/tern

Docs is at numeris.json. You can copy and test at [swagger](https://editor-next.swagger.io)

