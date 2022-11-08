package repositories

import (
	"db_forum/app/models"
	"db_forum/pkg/handlerows"
	"db_forum/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetInfoAboutUser(nickname string) (*models.User, error)
	GetSimilarUsers(user *models.User) (*[]models.User, error)
}

type UserRepositoryImpl struct {
	db *pgx.ConnPool
}

func MakeUserRepository(db *pgx.ConnPool) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (userRepository *UserRepositoryImpl) CreateUser(user *models.User) error {
	_, err := userRepository.db.Exec(queries.UserCreate, user.Nickname, user.Fullname, user.About, user.Email)
	return err
}

func (userRepository *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return userRepository.db.QueryRow(queries.UserUpdate, user.Fullname, user.About, user.Email, user.Nickname).Scan(&user.Fullname, &user.About, &user.Email)
}

func (userRepository *UserRepositoryImpl) GetInfoAboutUser(nickname string) (*models.User, error) {
	user := new(models.User)
	err := userRepository.db.QueryRow(queries.UserGet, nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return user, err
}

func (userRepository *UserRepositoryImpl) GetSimilarUsers(user *models.User) (*[]models.User, error) {
	resultRows, err := userRepository.db.Query(queries.UserGetSimilar, user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}
	defer resultRows.Close()
	return handlerows.User(resultRows)
}
