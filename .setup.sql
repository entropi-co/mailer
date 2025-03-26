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

insert into inbounds (id, content, sender, delivered_at)
values (1, 'content', 'sender@sender.dev', now()),
       (2, 'content2', 'sender@sender.dev', now()),
       (3, 'content3', 'sender@sender.dev', now()),
       (4, 'content4', 'sender@sender.dev', now())
on conflict do nothing;

insert into mailboxes (id, name, display_name, owner, type)
values (0, 'mailbox', 'mailbox', 0, 'inbox')
on conflict do nothing;

insert into mailboxes (id, name, display_name, owner, type)
values (1, 'mailbox2', 'mailbox2', 0, 'inbox')
on conflict do nothing;

insert into inbounds_mailboxes (inbound, mailbox, uid)
values (generate_series(1, 4), 0, generate_series(1, 4))
on conflict do nothing;

insert into inbounds_mailboxes (inbound, mailbox, uid)
values (2, 1, 0)
on conflict do nothing;

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