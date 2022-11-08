package handlers

import (
	"db_forum/app/models"
	"db_forum/app/usecases"
	"db_forum/pkg"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecases.UserUsecase
}

func MakeUserHandler(userUsecase_ usecases.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase_}
}

func (userHandler *UserHandler) CreateUser(c *gin.Context) {
	nickname := c.Param("nickname")

	var user models.User
	err := easyjson.UnmarshalFromReader(c.Request.Body, &user)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
		return
	}
	user.Nickname = nickname
	users, err := userHandler.userUsecase.CreateNewUser(&user)
	if err != nil && pkg.ConvertErrorToCode(err) != http.StatusConflict {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	if err != nil && pkg.ConvertErrorToCode(err) == http.StatusConflict {
		usersJson, internalErr := users.MarshalJSON()
		if internalErr != nil {
			c.Data(pkg.CreateErrorResponse(err))
			return
		}
		c.Data(pkg.ConvertErrorToCode(err), "application/json; charset=utf-8", usersJson)
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", userJSON)
}

func (userHandler *UserHandler) GetUser(c *gin.Context) {
	nickname := c.Param("nickname")

	user, err := userHandler.userUsecase.GetInfoAboutUser(nickname)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
}

func (userHandler *UserHandler) UpdateUser(c *gin.Context) {
	nickname := c.Param("nickname")

	userUpdate := new(models.UserUpdate)
	err := easyjson.UnmarshalFromReader(c.Request.Body, userUpdate)
	if err != nil {
		user, err := userHandler.userUsecase.GetInfoAboutUser(nickname)
		if err != nil {
			c.Data(pkg.CreateErrorResponse(err))
			return
		}

		userJSON, err := user.MarshalJSON()
		if err != nil {
			c.Data(pkg.CreateErrorResponse(err))
			return
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
		return
	}

	user := &models.User{Nickname: nickname, Fullname: userUpdate.Fullname, About: userUpdate.About, Email: userUpdate.Email}

	err = userHandler.userUsecase.UpdateUser(user)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
}
