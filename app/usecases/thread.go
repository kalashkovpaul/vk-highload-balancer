package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
	"strconv"
)

type ThreadUsecase interface {
	CreateNewPosts(slugOrID string, posts *models.Posts) error
	GetInfoAboutThread(slugOrID string) (*models.Thread, error)
	UpdateThread(slugOrID string, thread *models.Thread) error
	GetThreadPosts(slugOrID string, limit, since int, sort string, desc bool) (*models.Posts, error)
	VoteForThread(slugOrID string, vote *models.Vote) (*models.Thread, error)
}

type ThreadUsecaseImpl struct {
	repoVote   repositories.VoteRepository
	repoThread repositories.ThreadRepository
	repoUser   repositories.UserRepository
	repoPost   repositories.PostRepository
}

func MakeThreadUseCase(vote repositories.VoteRepository, thread repositories.ThreadRepository,
	user repositories.UserRepository, post repositories.PostRepository) ThreadUsecase {
	return &ThreadUsecaseImpl{repoVote: vote, repoThread: thread, repoUser: user, repoPost: post}
}

func (threadUsecase *ThreadUsecaseImpl) CreateNewPosts(slugOrID string, posts *models.Posts) error {
	var thread *models.Thread
	var err error
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	if err != nil {
		return pkg.ErrThreadNotFound
	}

	if len(*posts) == 0 {
		return err
	}

	if (*posts)[0].Parent != 0 {
		var parentPost *models.Post
		parentPost, err = threadUsecase.repoPost.GetPost((*posts)[0].Parent)
		if parentPost.Thread != thread.Id {
			return pkg.ErrParentPostFromOtherThread
		}
	}
	_, err = threadUsecase.repoUser.GetInfoAboutUser((*posts)[0].Author)
	if err != nil {
		return pkg.ErrUserNotFound
	}

	err = threadUsecase.repoThread.CreateThreadPosts(thread, posts)
	return err
}

func (threadUsecase *ThreadUsecaseImpl) GetInfoAboutThread(slugOrID string) (*models.Thread, error) {
	var thread *models.Thread
	var err error
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}
	if err != nil {
		return nil, pkg.ErrThreadNotFound
	}
	return thread, err
}

func (threadUsecase *ThreadUsecaseImpl) UpdateThread(slugOrID string, thread *models.Thread) error {
	id, errConv := strconv.Atoi(slugOrID)
	var currentThread *models.Thread
	var err error
	if errConv != nil {
		currentThread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		currentThread, err = threadUsecase.repoThread.GetById(int64(id))
	}
	if err != nil {
		return pkg.ErrThreadNotFound
	}
	if thread.Title != "" {
		currentThread.Title = thread.Title
	}
	if thread.Message != "" {
		currentThread.Message = thread.Message
	}
	err = threadUsecase.repoThread.UpdateThread(currentThread)
	if err != nil {
		return err
	}
	*thread = *currentThread
	return err
}

func (threadUsecase *ThreadUsecaseImpl) GetThreadPosts(slugOrID string, limit, since int, sort string, desc bool) (*models.Posts, error) {
	id, errConv := strconv.Atoi(slugOrID)
	var thread *models.Thread
	var err error
	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	if err != nil {
		return nil, pkg.ErrThreadNotFound
	}

	postsSlice := new([]models.Post)
	switch sort {
	case "tree":
		postsSlice, err = threadUsecase.repoThread.GetThreadPostsTree(thread.Id, limit, since, desc)
	case "parent_tree":
		postsSlice, err = threadUsecase.repoThread.GetThreadPostsParentTree(thread.Id, limit, since, desc)
	default:
		postsSlice, err = threadUsecase.repoThread.GetThreadPostsFlat(thread.Id, limit, since, desc)
	}
	if err != nil {
		return nil, err
	}
	posts := new(models.Posts)
	if len(*postsSlice) == 0 {
		*posts = []models.Post{}
	} else {
		*posts = *postsSlice
	}
	return posts, nil
}

func (threadUsecase *ThreadUsecaseImpl) VoteForThread(slugOrID string, vote *models.Vote) (*models.Thread, error) {
	var thread *models.Thread
	var err error
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	err = threadUsecase.repoVote.VoteForThread(thread.Id, vote)
	if err != nil {
		return nil, pkg.ErrUserNotFound
	}
	thread.Votes, err = threadUsecase.repoThread.GetThreadVotes(thread.Id)
	return thread, err
}
