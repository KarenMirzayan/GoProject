# Messenger
 ## Description
This is a REST api for messenger app. You can here use CRUD on users, conversations and messages.

## How to run?
### 1. You can start app through cli passing flags, specified in main.go
```
migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
port       = fs.Int("port", 8081, "API server port")
env        = fs.String("env", "development", "Environment (development|staging|production)")
dbDsn      = fs.String("dsn", "postgres://beezy:2202264mir@localhost:5433/messenger?sslmode=disable", "PostgreSQL DSN")
```
#### Example:
```
go run ./cmd/messenger \
-dsn="postgres://password:pa55word@localhost:5432/messenger?sslmode=disable" \
-migrations=file://pkg/messenger/migrations \
-env=development \
-port=8081
```

#### List of flags
`dsn` — postgress connection string with username, password, address, port, database name, and SSL mode. Default: `Value is not correct by security reasons`.

`migrations` — Path to folder with migration files. If not provided, migrations do not applied.

`env` - App running mode. Default: `development`

`port` - App port. Default: `8081`


### 2. You can build and run docker container with passing variables from .env
#### Example:
```
docker-compose --env-file .env.example up --build
```

❗IMPORTANT: Host value in DSN must have name of service from docker-compose. In our case hostname is `db`. 

Also, port have to be right side value after semicolon, it's a port of service available in docker isolated network, the left port for access from outside (host OS). In our case:
`"5433:5432"` port `5432` for docker isolated network.

Overall, your DSN for docker should be like this:
`postgres://postgres:postgres@db:5432/example?sslmode=disable`.

`--build` flag force docker compose to rebuild app. For example, if you have changed source code, you need this flag.

## CRUD
/users method POST

/users/{userId:[0-9]+} method GET

/users/{userId:[0-9]+} method PUT

/users/{userId:[0-9]+} method DELETE

## Postgres DB structure

```
Table users {
    user_id serial [primary key]
    firstname text [not null]
    lastname text [not null]
    date_of_birth date [not null]
    login varchar(16) [not null]
    password varchar(16) [not null]
}

Table user_conversations {
    conversation_id serial [primary key]
    user_id int [ref: <> users.user_id]
    friend_id int [ref: <> users.user_id]
}

Table messages {
    message_id serial [primary key]
    conversation_id int [ref: > user_conversations.conversation_id]
    sender_id int [ref: > users.user_id]
    content text [not null]
    timestamp timestamp(0) [not null, default: now()]
}
```
