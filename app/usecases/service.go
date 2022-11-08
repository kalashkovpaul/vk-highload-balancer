package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
)

type ServiceUsecase interface {
	ClearService() error
	GetService() (*models.Status, error)
}

type ServiceUsecaseImpl struct {
	repoService repositories.ServiceRepository
}

func MakeServiceUseCase(service repositories.ServiceRepository) ServiceUsecase {
	return &ServiceUsecaseImpl{repoService: service}
}

func (serviceUsecase *ServiceUsecaseImpl) ClearService() error {
	return serviceUsecase.repoService.ClearService()
}

func (serviceUsecase *ServiceUsecaseImpl) GetService() (*models.Status, error) {
	return serviceUsecase.repoService.GetService()
}
