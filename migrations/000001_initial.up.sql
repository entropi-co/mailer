CREATE TABLE users
(
    id         int8 primary key,
    local      varchar(64)              not null unique,
    created_at timestamp with time zone not null default now()
);

CREATE TABLE keys
(
    value varchar(128) not null unique,
    owner int8         null references users (id)
);

CREATE TYPE mailbox_type AS ENUM ('inbox', 'junk', 'sent', 'drafts', 'trash', 'user');

CREATE TABLE mailboxes
(
    id           int8 primary key,
    name         varchar(255)             not null,
    display_name varchar(255)             null,
    owner        int8                     null references users (id),
    type         mailbox_type             not null default 'user',
    created_at   timestamp with time zone not null default now(),
    unique (owner, name)
);

CREATE TABLE inbounds
(
    id           int8 primary key,
    header       jsonb not null default '{}',
    content      text,
    sender       varchar(320)             not null,
    delivered_at timestamp with time zone not null default now()
);

CREATE TABLE inbounds_mailboxes
(
    inbound int8 references inbounds (id),
    mailbox int8 references mailboxes (id),
    uid     int4 not null,
    unique (inbound, mailbox),
    unique (uid, mailbox)
);