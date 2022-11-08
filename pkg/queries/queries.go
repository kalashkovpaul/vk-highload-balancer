package queries

var (
	ForumCreate              = `insert into "forums" ("title", "user_", "slug") values ($1, $2, $3);`
	ForumGetBySlug           = `select "title", "user_", "slug", "posts", "threads" from "forums" where "slug" = $1`
	ForumGetUsers            = "select users.nickname, users.fullname, users.about, users.email from users left join user_forum on users.nickname = user_forum.nickname where user_forum.forum = $1 order by users.nickname limit $2;"
	ForumGetUsersDesc        = "select users.nickname, users.fullname, users.about, users.email from users left join user_forum on users.nickname = user_forum.nickname where user_forum.forum = $1 order by users.nickname desc limit $2;"
	ForumGetUsersSince       = "select users.nickname, users.fullname, users.about, users.email from users left join user_forum on users.nickname = user_forum.nickname where user_forum.forum = $1 and users.nickname > $2 order by users.nickname limit $3;"
	ForumGetUsersSinceDesc   = "select users.nickname, users.fullname, users.about, users.email from users left join user_forum on users.nickname = user_forum.nickname where user_forum.forum = $1 and users.nickname < $2 order by users.nickname desc limit $3;"
	ForumGetThreads          = "select id, title, author, forum, message, votes, slug, created from threads where forum = $1 order by created asc limit $2;"
	ForumGetThreadsDesc      = "select id, title, author, forum, message, votes, slug, created from threads where forum = $1 order by created desc limit $2;"
	ForumGetThreadsSince     = "select id, title, author, forum, message, votes, slug, created from threads where forum = $1 and created >= $2 order by created asc limit $3;"
	ForumGetThreadsSinceDesc = "select id, title, author, forum, message, votes, slug, created from threads where forum = $1 and created <= $2 order by created desc limit $3;"

	PostGet    = "select id, coalesce(parent, 0), author, message, is_edited, forum, thread, created from posts where id = $1"
	PostUpdate = "update posts set message = $1, is_edited = $2 where id = $3;"
	PostPart   = "insert into posts (parent, author, message, forum, thread, created) values "

	ServiceClear = "truncate table forums, posts, threads, user_forum, users, votes;"
	ServiceGet   = "select (select count(*) from users) as users, (select count(*) from forums) as forums, (select count(*) from threads) as threads, (select count(*) from posts) as posts;"

	ThreadCreate  = "insert into threads (title, author, forum, message, slug, created) values ($1, $2, $3, $4, $5, $6) returning id, created;"
	ThreadGetSlug = "select id, title, author, forum, message, votes, slug, created from threads where slug = $1;"
	ThreadGetId   = "select id, title, author, forum, message, votes, slug, created from threads where id = $1;"
	ThreadVotes   = "select votes from threads where id = $1;"
	ThreadUpdate  = "update threads SET title = $1, message = $2 where id = $3;"

	ThreadFlatBase      = "select id, coalesce(parent, 0), author, message, is_edited, forum, thread, created from posts where thread = $1 "
	ThreadFlat          = "and id > $2 order by id limit $3;"
	ThreadFlatDesc      = "and id < $2 order by id desc limit $3;"
	ThreadFlatSince     = "order by id limit $2;"
	ThreadFlatSinceDesc = " order by id desc limit $2;"

	ThreadTreeBase      = "select id, coalesce(parent, 0), author, message, is_edited, forum, thread, created from posts "
	ThreadTree          = "where thread = $1 and path > (select path from posts where id = $2) order by path limit $3;"
	ThreadTreeDesc      = "where thread = $1 and path < (select path from posts where id = $2) order by path desc limit $3;"
	ThreadTreeSince     = "where thread = $1 order by path limit $2;"
	ThreadTreeSinceDesc = "where thread = $1 order by path desc limit $2;"

	ThreadParentBase          = "select id, coalesce(parent, 0), author, message, is_edited, forum, thread, created from posts where path[1] in "
	ThreadParentTree          = "(select id from posts where thread = $1 and parent is null and path[1] > (select path[1] from posts where id = $2) order by path[1] limit $3) order by path;"
	ThreadParentTreeDesc      = "(select id from posts where thread = $1 and parent is null and path[1] < (select path[1] from posts where id = $2) order by path[1] desc limit $3) order by path[1] desc, path [2:];"
	ThreadParentTreeSince     = "(select id from posts where thread = $1 and parent is null order by path[1] limit $2) order by path;"
	ThreadParentTreeSinceDesc = "(select id from posts where thread = $1 and parent is null order by path[1] desc limit $2) order by path[1] desc, path[2:]"

	UserCreate     = "insert into users values ($1, $2, $3, $4);"
	UserUpdate     = "update users set fullname = $1, about = $2, email = $3 where nickname = $4 returning fullname, about, email;"
	UserGet        = "select nickname, fullname, about, email from users where nickname = $1;"
	UserGetSimilar = "select nickname, fullname, about, email from users where nickname = $1 or email = $2;"

	Vote = "insert into votes (nickname, thread, voice) values ($1, $2, $3) on conflict (nickname, thread) do update set voice = excluded.voice;"
)
