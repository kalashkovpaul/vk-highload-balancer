package handlers

import (
	"db_forum/app/models"
	"db_forum/app/usecases"
	"db_forum/pkg"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type ThreadHandler struct {
	threadUsecase usecases.ThreadUsecase
}

func MakeThreadHandler(threadUsecase_ usecases.ThreadUsecase) *ThreadHandler {
	return &ThreadHandler{threadUsecase: threadUsecase_}
}

func (threadHandler *ThreadHandler) CreatePosts(c *gin.Context) {
	rawId := c.Param("slug_or_id")
	var posts models.Posts
	err := easyjson.UnmarshalFromReader(c.Request.Body, &posts)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
		return
	}

	err = threadHandler.threadUsecase.CreateNewPosts(rawId, &posts)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	postsJSON, err := posts.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", postsJSON)
}

func (threadHandler *ThreadHandler) GetThread(c *gin.Context) {
	rawId := c.Param("slug_or_id")

	thread, err := threadHandler.threadUsecase.GetInfoAboutThread(rawId)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadJSON)
}

func (threadHandler *ThreadHandler) UpdateThread(c *gin.Context) {
	rawId := c.Param("slug_or_id")

	var threadUpdate models.ThreadUpdate
	err := easyjson.UnmarshalFromReader(c.Request.Body, &threadUpdate)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
		return
	}

	thread := &models.Thread{Title: threadUpdate.Title, Message: threadUpdate.Message}
	err = threadHandler.threadUsecase.UpdateThread(rawId, thread)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadJSON)
}

func (threadHandler *ThreadHandler) GetThreadPosts(c *gin.Context) {
	rawId := c.Param("slug_or_id")

	since := -1
	rawSince := c.Query("since")
	if rawSince != "" {
		var err error
		since, err = strconv.Atoi(rawSince)
		if err != nil {
			c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
			return
		}
	}

	rawLimit := c.Query("limit")
	defaultLimit := 100

	if rawLimit != "" {
		var err error
		defaultLimit, err = strconv.Atoi(rawLimit)
		if err != nil {
			c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
			return
		}
	}

	rawDecs := c.Query("desc")
	defaultDesc := false

	if rawDecs != "" {
		var err error
		defaultDesc, err = strconv.ParseBool(rawDecs)
		if err != nil {
			c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
			return
		}
	}

	sort := c.Query("sort")
	if sort == "" {
		sort = "flat"
	}

	posts, err := threadHandler.threadUsecase.GetThreadPosts(rawId, defaultLimit, since, sort, defaultDesc)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	postsJSON, err := posts.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", postsJSON)
}

func (threadHandler *ThreadHandler) Vote(c *gin.Context) {
	rawId := c.Param("slug_or_id")

	var vote models.Vote
	err := easyjson.UnmarshalFromReader(c.Request.Body, &vote)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
		return
	}

	thread, err := threadHandler.threadUsecase.VoteForThread(rawId, &vote)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadJSON)
}
