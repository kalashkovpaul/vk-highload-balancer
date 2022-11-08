package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (err error)
	GetInfoAboutForum(slug string) (forum *models.Forum, err error)
	CreateForumsThread(thread *models.Thread) (err error)
	GetForumUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error)
	GetForumThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error)
}

type ForumUseCaseImpl struct {
	repoForum  repositories.ForumRepository
	repoThread repositories.ThreadRepository
	repoUser   repositories.UserRepository
}

func MakeForumUseCase(forum repositories.ForumRepository, thread repositories.ThreadRepository, user repositories.UserRepository) *ForumUseCaseImpl {
	return &ForumUseCaseImpl{repoForum: forum, repoThread: thread, repoUser: user}
}

func (forumUsecase *ForumUseCaseImpl) CreateForum(forum *models.Forum) error {
	user, err := forumUsecase.repoUser.GetInfoAboutUser(forum.User)
	if err != nil {
		return pkg.ErrUserNotFound
	}

	oldForum, err := forumUsecase.repoForum.GetInfoAboutForum(forum.Slug)
	if oldForum.Slug != "" {
		*forum = *oldForum
		return pkg.ErrForumAlreadyExists
	}

	forum.User = user.Nickname
	err = forumUsecase.repoForum.CreateForum(forum)
	return err
}

func (forumUsecase *ForumUseCaseImpl) GetInfoAboutForum(slug string) (*models.Forum, error) {
	forum, err := forumUsecase.repoForum.GetInfoAboutForum(slug)
	if err != nil {
		return nil, pkg.ErrForumNotExist
	}
	return forum, nil
}

func (forumUsecase *ForumUseCaseImpl) CreateForumsThread(thread *models.Thread) error {
	forum, err := forumUsecase.repoForum.GetInfoAboutForum(thread.Forum)
	if err != nil {
		return pkg.ErrForumOrTheadNotFound
	}

	_, err = forumUsecase.repoUser.GetInfoAboutUser(thread.Author)
	if err != nil {
		return pkg.ErrForumOrTheadNotFound
	}

	currentThread, err := forumUsecase.repoThread.GetBySlug(thread.Slug)
	if currentThread.Slug != "" {
		*thread = *currentThread
		return pkg.ErrThreadAlreadyExists
	}

	thread.Forum = forum.Slug
	err = forumUsecase.repoThread.CreateThread(thread)
	return err
}

func (forumUsecase *ForumUseCaseImpl) GetForumUsers(slug string, limit int, since string, desc bool) (*models.Users, error) {
	_, err := forumUsecase.repoForum.GetInfoAboutForum(slug)
	if err != nil {
		return nil, pkg.ErrForumNotExist
	}

	usersSlice, err := forumUsecase.repoForum.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}
	users := new(models.Users)
	if len(*usersSlice) == 0 {
		*users = []models.User{}
	} else {
		*users = *usersSlice
	}
	return users, err
}

func (forumUsecase *ForumUseCaseImpl) GetForumThreads(slug string, limit int, since string, desc bool) (*models.Threads, error) {
	forum, err := forumUsecase.repoForum.GetInfoAboutForum(slug)
	if err != nil {
		return nil, pkg.ErrForumNotExist
	}

	threadsSlice, err := forumUsecase.repoForum.GetForumThreads(forum.Slug, limit, since, desc)
	if err != nil {
		return nil, err
	}
	threads := new(models.Threads)
	if len(*threadsSlice) == 0 {
		*threads = []models.Thread{}
	} else {
		*threads = *threadsSlice
	}
	return threads, err
}
