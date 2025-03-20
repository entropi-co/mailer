CREATE TABLE users
(
    id         bigint primary key,
    local      varchar(64)              not null unique,
    created_at timestamp with time zone not null default now()
);

CREATE TABLE keys
(
    value varchar(128) not null unique,
    owner bigint       null references users (id)
);

CREATE TYPE mailbox_type AS ENUM ('inbox', 'junk', 'sent', 'drafts', 'trash', 'user');

CREATE TABLE mailboxes
(
    id           bigint primary key,
    name         varchar(255)             not null,
    display_name varchar(255)             null,
    owner        bigint                   null references users (id),
    type         mailbox_type             not null default 'user',
    created_at   timestamp with time zone not null default now()
);

CREATE TABLE inbounds
(
    id           bigint primary key,
    content      text,
    sender       varchar(320)             not null,
    delivered_at timestamp with time zone not null default now()
);

CREATE TABLE inbounds_mailboxes
(
    inbound bigint references inbounds (id),
    mailbox bigint references mailboxes (id),
    unique (inbound, mailbox)
);