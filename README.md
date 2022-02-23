# Taskeroo

Task management, primarily meant for managing tasks between residents living together.

URL to production: https://taskeroo.tychsen.me

## Start locally

First, you need a database (PostgreSQL):

Run docker-compose, which will start a local instance of postgres:
```
docker-compose up -d
```

Run the application, either using an IDE or the command line:
```shell
go run main.go
```

If you need to test Telegram, you will need to specify a Telegram token, either in your IDE
or before you execute `go run`:
```shell
TELEGRAM_TOKEN=bla go run main.go
```
or
```shell
export TELEGRAM_TOKEN=bla
go run main.go
```
