CREATE TABLE users
(
    id         bigint primary key,
    local      varchar(64)              not null unique,
    created_at timestamp with time zone not null default now()
);

CREATE TABLE inbounds
(
    id           bigint primary key,
    content      text,
    sender       varchar(320)             not null,
    delivered_at timestamp with time zone not null default now()
);

CREATE TABLE inbounds_recipients
(
    inbound   bigint references inbounds (id),
    recipient bigint references users (id),
    unique (inbound, recipient)
);