package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
	"strings"
)

type PostUsecase interface {
	GetInfoAboutPost(id int64, related string) (*models.PostFull, error)
	UpdatePost(post *models.Post) (err error)
}

type PostUsecaseImpl struct {
	repoForum  repositories.ForumRepository
	repoThread repositories.ThreadRepository
	repoUser   repositories.UserRepository
	repoPost   repositories.PostRepository
}

func MakePostUseCase(forum repositories.ForumRepository, thread repositories.ThreadRepository,
	user repositories.UserRepository, post repositories.PostRepository) PostUsecase {
	return &PostUsecaseImpl{repoForum: forum, repoThread: thread, repoUser: user, repoPost: post}
}

func (postUsecase *PostUsecaseImpl) GetInfoAboutPost(id int64, related string) (*models.PostFull, error) {
	fullPost := new(models.PostFull)
	var post *models.Post
	var err error
	post, err = postUsecase.repoPost.GetPost(id)
	if err != nil {
		err = pkg.ErrPostNotFound
	}
	fullPost.Post = post

	var relatedDataArr []string
	if related != "" {
		relatedDataArr = strings.Split(related, ",")
	}

	for _, data := range relatedDataArr {
		switch data {
		case "thread":
			var thread *models.Thread
			thread, err = postUsecase.repoThread.GetById(fullPost.Post.Thread)
			if err != nil {
				err = pkg.ErrThreadNotFound
			}
			fullPost.Thread = thread
		case "user":
			var user *models.User
			user, err = postUsecase.repoUser.GetInfoAboutUser(fullPost.Post.Author)
			if err != nil {
				err = pkg.ErrUserNotFound
			}
			fullPost.Author = user
		case "forum":
			var forum *models.Forum
			forum, err = postUsecase.repoForum.GetInfoAboutForum(fullPost.Post.Forum)
			if err != nil {
				err = pkg.ErrForumNotExist
			}
			fullPost.Forum = forum
		}
	}
	return fullPost, err
}

func (postUsecase *PostUsecaseImpl) UpdatePost(post *models.Post) error {
	currentPost, err := postUsecase.repoPost.GetPost(post.Id)
	if err != nil {
		return pkg.ErrThreadNotFound
	}

	if post.Message != "" {
		if currentPost.Message != post.Message {
			currentPost.IsEdited = true
		}
		currentPost.Message = post.Message
		err = postUsecase.repoPost.UpdatePost(currentPost)
		if err != nil {
			return err
		}
	}
	*post = *currentPost
	return nil
}
