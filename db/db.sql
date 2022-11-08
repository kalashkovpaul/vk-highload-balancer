create extension if not exists citext;

create unlogged table if not exists users
(
    nickname citext collate "C" not null primary key,
    fullname citext             not null,
    about    text,
    email    citext             not null unique
);

create unlogged table if not exists forums
(
    title   text   not null,
    user_   citext not null references users (nickname) on update cascade on delete cascade,
    slug    citext not null primary key,
    posts   int default 0,
    threads int default 0
);

create unlogged table if not exists threads
(
    id      serial not null primary key,
    title   text   not null,
    author  citext not null references users (nickname) on update cascade on delete cascade,
    forum   citext not null references forums (slug) on update cascade on delete cascade,
    message text   not null,
    votes   int                      default 0,
    slug    citext,
    created timestamp with time zone default now()
);

create unlogged table if not exists posts
(
    id        bigserial not null primary key,
    parent    int references posts (id),
    author    citext not null references users (nickname),
    message   text   not null,
    is_edited bool                     default false,
    forum     citext not null references forums (slug),
    thread    int    not null references threads (id),
    created   timestamp with time zone default now(),
    path      bigint[]                 default array []::integer[]
);

create unlogged table if not exists votes
(
    nickname citext not null references users (nickname),
    thread   int    not null references threads (id),
    voice    int    not null,
    constraint user_thread_key unique (nickname, thread)
);

create unlogged table if not exists user_forum
(
    nickname citext collate "C" not null references users (nickname),
    forum    citext             not null references forums (slug),
    constraint user_forum_key unique (nickname, forum)
);

-- Триггеры и процедуры
create or replace function create_user()
    returns trigger as
$$
begin
    insert into user_forum (nickname, forum) values (new.author, new.forum) on conflict do nothing;
    return new;
end;
$$ language plpgsql;

create trigger create_new_thread
    after insert
    on threads
    for each row
execute procedure create_user();

create trigger create_new_post
    after insert
    ON posts
    for each row
execute procedure create_user();

create or replace function create_post_before()
    returns trigger as
$$
begin
    new.path = (select path from posts where id = new.parent) || new.id;
    return new;
end;
$$ language plpgsql;

create trigger create_post_before
    before insert
    on posts
    for each row
execute procedure create_post_before();

create or replace function create_post_after()
    returns trigger as
$$
begin
    update forums set posts = forums.posts + 1 where slug = new.forum;
    return new;
end;
$$ language plpgsql;

create trigger create_post_after
    after insert
    on posts
    for each row
execute procedure create_post_after();

create or replace function create_votes()
    returns trigger as
$$
begin
    update threads set votes = votes + new.voice where id = new.thread;
    return new;
end;
$$ language plpgsql;

create trigger create_votes
    after insert
    on votes
    for each row
execute procedure create_votes();


create or replace function update_votes()
    returns trigger as
$$
begin
    update threads set votes = votes - old.voice + new.voice where id = new.thread;
    return NULL;
end;
$$ language plpgsql;

create trigger update_votes
    after update
    on votes
    for each row
execute procedure update_votes();

create or replace function create_thread()
    returns trigger as
$$
begin
    update forums set threads = forums.threads + 1 where slug = new.forum;
    RETURN new;
end;
$$ language plpgsql;

create trigger create_thread
    after insert
    on threads
    for each row
execute procedure create_thread();

CREATE INDEX IF NOT EXISTS users_idx on users (nickname, email) include (about, fullname);
create index if not exists users_nickname_hash on users using hash (nickname);
create index if not exists user_forum_all on user_forum (forum, nickname);

create index if not exists forums_slug on forums using hash (slug);

create index if not exists threads_created on threads using hash (created);
create index if not exists threads_slug on threads using hash (slug);
create index if not exists threads_forum ON threads using hash (forum);
create index if not exists threads_forum_created on threads (forum, created);
create index if not exists threads_id ON threads USING hash (id);

create index if not exists posts_id on posts using hash (id);
create index if not exists posts_threads on posts using hash (thread);
create index if not exists posts_threads_parent on posts ((path[1]), path);
create index if not exists posts_threads_past on posts (thread, path);
create index if not exists posts_threads_id ON posts (thread, id);
create index if not exists posts_threads_path ON posts (thread, (path[1]));

create unique index if not exists votes_nickname on votes (thread, nickname);