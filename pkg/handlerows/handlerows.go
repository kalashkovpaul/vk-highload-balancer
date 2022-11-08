package handlerows

import (
	"db_forum/app/models"
	"github.com/jackc/pgx"
	"time"
)

func Thread(result *pgx.Rows) (*[]models.Thread, error) {
	var err error
	var bufThreads []models.Thread
	for result.Next() {
		thread := models.Thread{}
		err = result.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)
		if err != nil {
			return nil, err
		}
		bufThreads = append(bufThreads, thread)
	}
	return &bufThreads, nil
}

func Post(result *pgx.Rows) (*[]models.Post, error) {
	posts := new([]models.Post)
	var err error
	for result.Next() {
		post := models.Post{}
		postTime := time.Time{}

		err = result.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &postTime)
		if err != nil {
			return nil, err
		}

		post.Created = postTime.Format(time.RFC3339)
		*posts = append(*posts, post)
	}
	return posts, nil
}

func User(result *pgx.Rows) (*[]models.User, error) {
	var users []models.User
	var err error
	for result.Next() {
		user := models.User{}
		err = result.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}
