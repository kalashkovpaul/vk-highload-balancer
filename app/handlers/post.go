package handlers

import (
	"db_forum/app/models"
	"db_forum/app/usecases"
	"db_forum/pkg"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUsecase usecases.PostUsecase
}

func MakePostHandler(postUsecase_ usecases.PostUsecase) *PostHandler {
	return &PostHandler{postUsecase: postUsecase_}
}

func (postHandler *PostHandler) GetPost(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.Atoi(rawId)

	related := c.Query("related")

	postFull, err := postHandler.postUsecase.GetInfoAboutPost(int64(id), related)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	postFullJSON, err := postFull.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", postFullJSON)
}

func (postHandler *PostHandler) UpdatePost(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	var postUpdate models.PostUpdate
	err = easyjson.UnmarshalFromReader(c.Request.Body, &postUpdate)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
		return
	}

	post := &models.Post{Id: int64(id), Message: postUpdate.Message}
	err = postHandler.postUsecase.UpdatePost(post)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	postJSON, err := post.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", postJSON)
}
