package handlers

import (
	"db_forum/app/usecases"
	"db_forum/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceUsecase usecases.ServiceUsecase
}

func MakeServiceHandler(serviceUsecase_ usecases.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{serviceUsecase: serviceUsecase_}
}

func (serviceHandler *ServiceHandler) Clear(c *gin.Context) {
	err := serviceHandler.serviceUsecase.ClearService()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}

func (serviceHandler *ServiceHandler) GetStatus(c *gin.Context) {
	status, err := serviceHandler.serviceUsecase.GetService()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	statusJSON, err := status.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", statusJSON)
}
