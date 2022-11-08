package repositories

import (
	"db_forum/app/models"
	"db_forum/pkg/handlerows"
	"db_forum/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type ForumRepository interface {
	CreateForum(forum *models.Forum) (err error)
	GetInfoAboutForum(slug string) (forum *models.Forum, err error)
	GetForumUsers(slug string, limit int, since string, desc bool) (*[]models.User, error)
	GetForumThreads(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error)
}

type ForumRepositoryImpl struct {
	db *pgx.ConnPool
}

func MakeForumRepository(db *pgx.ConnPool) ForumRepository {
	return &ForumRepositoryImpl{db: db}
}

func (forumRepository *ForumRepositoryImpl) CreateForum(forum *models.Forum) (err error) {
	_, err = forumRepository.db.Exec(queries.ForumCreate, forum.Title, forum.User, forum.Slug)
	return err
}

func (forumRepository *ForumRepositoryImpl) GetInfoAboutForum(slug string) (forum *models.Forum, err error) {
	forum = new(models.Forum)
	err = forumRepository.db.QueryRow(queries.ForumGetBySlug, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return forum, err
}

func (forumRepository *ForumRepositoryImpl) GetForumUsers(slug string, limit int, since string, desc bool) (*[]models.User, error) {
	var query string

	var result *pgx.Rows
	var innerError error

	if since != "" {
		if desc {
			query = queries.ForumGetUsersSinceDesc
		} else {
			query = queries.ForumGetUsersSince
		}
		result, innerError = forumRepository.db.Query(query, slug, since, limit)
		if innerError != nil {
			return nil, innerError
		}
	} else {
		if desc {
			query = queries.ForumGetUsersDesc
		} else {
			query = queries.ForumGetUsers
		}
		result, innerError = forumRepository.db.Query(query, slug, limit)
		if innerError != nil {
			return nil, innerError
		}
	}
	defer result.Close()
	return handlerows.User(result)
}

func (forumRepository *ForumRepositoryImpl) GetForumThreads(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error) {
	var query string

	var result *pgx.Rows
	var innerError error

	if since != "" {
		if desc {
			query = queries.ForumGetThreadsSinceDesc
		} else {
			query = queries.ForumGetThreadsSince
		}
		result, innerError = forumRepository.db.Query(query, slug, since, limit)
		if innerError != nil {
			return
		}
	} else {
		if desc {
			query = queries.ForumGetThreadsDesc
		} else {
			query = queries.ForumGetThreads
		}
		result, innerError = forumRepository.db.Query(query, slug, limit)
		if innerError != nil {
			return
		}
	}

	defer result.Close()
	return handlerows.Thread(result)
}
