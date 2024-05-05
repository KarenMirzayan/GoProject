create extension if not exists citext;

CREATE TABLE IF NOT EXISTS users
(
    id            bigserial PRIMARY KEY,
    created_at    timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name          text                        NOT NULL,
    email         citext UNIQUE               NOT NULL,
    password_hash bytea                       NOT NULL,
    activated     bool                        NOT NULL,
    version       integer                     NOT NULL DEFAULT 1
);

create table if not exists user_conversations
(
    conversation_id serial primary key,
    user_id         int,
    friend_id       int,
    foreign key (user_id) references users (id) on delete cascade,
    foreign key (friend_id) references users (id) on delete cascade,
    constraint unique_conversation unique (user_id, friend_id)
);

create table if not exists messages
(
    message_id      serial primary key,
    conversation_id int,
    sender_id       int,
    content         text         not null,
    timestamp       timestamp(0) not null default now(),
    foreign key (conversation_id) references user_conversations (conversation_id) on delete cascade,
    foreign key (sender_id) references users (id) on delete cascade
);

create table if not exists channels
(
    channel_id serial primary key,
    user_id    int,
    name       text,
    foreign key (user_id) references users (id) on delete cascade
);

CREATE TABLE IF NOT EXISTS tokens
(
    hash    BYTEA PRIMARY KEY,
    user_id BIGINT                      NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry  TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope   TEXT                        NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions
(
    id   BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions
(
    user_id       BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES ('conversation:read'),
       ('conversation:write');