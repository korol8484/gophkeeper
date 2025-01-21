create table if not exists "user"
(
    id            bigserial
        constraint user_pk
            primary key,
    login         varchar(100) not null,
    password_hash varchar(255) not null
);

create unique index if not exists user_login_uindex
    on "user" (login);
