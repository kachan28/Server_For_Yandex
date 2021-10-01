create database discontdealer;

\c discontdealer

create table if not exists users(
    id serial       primary key,
    username        varchar(15) not null,
    password        char(128)   not null,
    salt            char(32)    not null,
    advert_disabled boolean     default false,
    is_premium      boolean     default false,
    referal_code    char(5)     not null,
    referal_count   smallint    default 0,
    created_at      timestamp,
    updated_at      timestamp
);