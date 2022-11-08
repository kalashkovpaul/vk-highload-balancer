package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
)

type UserUsecase interface {
	CreateNewUser(user *models.User) (*models.Users, error)
	GetInfoAboutUser(nickname string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type UserUsecaseImpl struct {
	repoUser repositories.UserRepository
}

func MakeUserUseCase(user repositories.UserRepository) UserUsecase {
	return &UserUsecaseImpl{repoUser: user}
}

func (userUsecase *UserUsecaseImpl) CreateNewUser(user *models.User) (*models.Users, error) {
	var users *models.Users
	similarUsers, err := userUsecase.repoUser.GetSimilarUsers(user)
	if err != nil {
		return users, pkg.ErrUserAlreadyExist
	} else if len(*similarUsers) > 0 {
		users = new(models.Users)
		*users = *similarUsers
		return users, pkg.ErrUserAlreadyExist
	}
	err = userUsecase.repoUser.CreateUser(user)
	return users, err
}

func (userUsecase *UserUsecaseImpl) GetInfoAboutUser(nickname string) (*models.User, error) {
	user, err := userUsecase.repoUser.GetInfoAboutUser(nickname)
	if err != nil {
		return nil, pkg.ErrUserNotFound
	}
	return user, nil
}

func (userUsecase *UserUsecaseImpl) UpdateUser(user *models.User) error {
	oldUser, err := userUsecase.repoUser.GetInfoAboutUser(user.Nickname)
	if oldUser.Nickname == "" {
		return pkg.ErrUserNotFound
	}
	if oldUser.Fullname != user.Fullname && user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}
	if oldUser.About != user.About && user.About == "" {
		user.About = oldUser.About
	}
	if oldUser.Email != user.Email && user.Email == "" {
		user.Email = oldUser.Email
	}
	err = userUsecase.repoUser.UpdateUser(user)
	if err != nil {
		return pkg.ErrUserDataConflict
	}
	return nil
}
