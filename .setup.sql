insert into users (id, local)
values (0, 'user')
on conflict do nothing;

insert into users (id, local)
values (1, 'user2')
on conflict do nothing;

insert into keys (value, owner)
values ('@kv', 0)
on conflict do nothing;

insert into keys (value, owner)
values ('@kv2', 1)
on conflict do nothing;

insert into inbounds (id, body, sender, delivered_at)
values (1, 'content'::bytea, 'sender@sender.dev', now()),
       (2, 'content2'::bytea, 'sender@sender.dev', now()),
       (3, 'content3'::bytea, 'sender@sender.dev', now()),
       (4, 'content4'::bytea, 'sender@sender.dev', now()),
       (5, 'content5'::bytea, 'sender@sender.dev', now())
on conflict do nothing;

insert into mailboxes (id, name, display_name, owner, type, priority)
values (0, 'mailbox', 'mailbox', 0, 'inbox', 0),
       (1, 'mailbox2', 'mailbox2', 0, 'inbox', 1),
       (2, 'mailbox3', 'mailbox3', 0, 'inbox', 2),
       (3, 'mailbox4', 'mailbox4', 0, 'inbox', 3),
       (4, 'mailbox5', 'mailbox5', 0, 'inbox', 4)
on conflict do nothing;

-- insert into inbounds_mailboxes (inbound, mailbox, uid)
-- values (generate_series(1, 4), 0, generate_series(1, 4))
-- on conflict do nothing;

-- insert into inbounds_mailboxes (inbound, mailbox, uid)
-- values (2, 1, 0)
-- on conflict do nothing;

WITH ranked_rows AS (SELECT inbounds.*, ROW_NUMBER() OVER (ORDER BY inbounds.id) AS sequence
                     FROM inbounds
                              left join public.inbounds_mailboxes im on inbounds.id = im.inbound
                     where mailbox = 1)
SELECT *
FROM ranked_rows
WHERE sequence IN (1, 3);

SELECT *
FROM keys
         LEFT JOIN public.users u on u.id = keys.owner
WHERE local = 'user';

-- # Query: Insert inbound
WITH updated AS (
    UPDATE mailboxes AS m
        SET uid_next = m.uid_next + 1
        WHERE id IN (SELECT id
                     FROM (SELECT m.id,
                                  ROW_NUMBER() OVER (PARTITION BY owner ORDER BY priority, m.created_at) AS _row
                           FROM mailboxes m
                                    LEFT JOIN public.users
                                              on users.id = m.owner
                           WHERE users.local IN ('user', 'user2')) as _sub
                     WHERE _row = 1)
        RETURNING m.id, uid_next)
INSERT
INTO inbounds_mailboxes (inbound, mailbox, uid)
SELECT 1, updated.id, updated.uid_next
FROM updated;