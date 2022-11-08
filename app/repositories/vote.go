package repositories

import (
	"db_forum/app/models"
	"db_forum/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type VoteRepository interface {
	VoteForThread(id int64, vote *models.Vote) error
}

type VoteRepositoryImpl struct {
	db *pgx.ConnPool
}

func MakeVoteRepository(db *pgx.ConnPool) VoteRepository {
	return &VoteRepositoryImpl{db: db}
}

func (voteRepository *VoteRepositoryImpl) VoteForThread(id int64, vote *models.Vote) error {
	_, err := voteRepository.db.Exec(queries.Vote, vote.Nickname, id, vote.Voice)
	return err
}
