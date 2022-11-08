package repositories

import (
	"db_forum/app/models"
	"db_forum/pkg/queries"
	"time"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type PostRepository interface {
	GetPost(id int64) (post *models.Post, err error)
	UpdatePost(post *models.Post) (err error)
}

type PostRepositoryImpl struct {
	db *pgx.ConnPool
}

func MakePostRepository(db *pgx.ConnPool) PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (postStore *PostRepositoryImpl) GetPost(id int64) (post *models.Post, err error) {
	post = &models.Post{}
	timeScan := time.Time{}
	err = postStore.db.QueryRow(queries.PostGet, id).
		Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&timeScan)
	post.Created = timeScan.Format(time.RFC3339)
	return
}

func (postStore *PostRepositoryImpl) UpdatePost(post *models.Post) (err error) {
	_, err = postStore.db.Exec(queries.PostUpdate, post.Message, post.IsEdited, post.Id)
	return
}
