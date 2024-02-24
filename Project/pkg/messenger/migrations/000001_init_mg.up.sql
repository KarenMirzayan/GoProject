create table if not exists users (
                                     user_id serial primary key,
                                     firstname text not null,
                                     lastname text not null,
                                     date_of_birth timestamp(0) with time zone not null,
                                     login varchar(16) not null,
                                     password varchar(16) not null
);

create table if not exists user_conversations (
                                    conversation_id serial primary key,
                                    user_id int,
                                    friend_id int,
                                    foreign key (user_id) references users(user_id),
                                    foreign key (friend_id) references users(user_id)
);

create table if not exists messages (
                          message_id serial primary key,
                          conversation_id int,
                          sender_id int,
                          content text not null,
                          timestamp timestamp(0) not null default now(),
                          foreign key (conversation_id) references user_conversations(conversation_id),
                          foreign key (sender_id) references users(user_id)
);