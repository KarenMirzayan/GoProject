## Description
This is a REST api for messenger app. You can here use CRUD on users, conversations and messages.

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
