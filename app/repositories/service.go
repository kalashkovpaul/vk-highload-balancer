package repositories

import (
	"db_forum/app/models"
	"db_forum/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type ServiceRepository interface {
	ClearService() (err error)
	GetService() (status *models.Status, err error)
}

type ServiceRepositoryImpl struct {
	db *pgx.ConnPool
}

func MakeServiceRepository(db *pgx.ConnPool) ServiceRepository {
	return &ServiceRepositoryImpl{db: db}
}

func (serviceRepository *ServiceRepositoryImpl) ClearService() (err error) {
	_, err = serviceRepository.db.Exec(queries.ServiceClear)
	return
}

func (serviceRepository *ServiceRepositoryImpl) GetService() (status *models.Status, err error) {
	status = &models.Status{}
	err = serviceRepository.db.QueryRow(queries.ServiceGet).
		Scan(
			&status.User,
			&status.Forum,
			&status.Thread,
			&status.Post)
	return
}
