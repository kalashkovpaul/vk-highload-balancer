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

type ForumHandler struct {
	forumUsecase usecases.ForumUsecase
}

func MakeForumHandler(forumUsecase_ usecases.ForumUsecase) *ForumHandler {
	return &ForumHandler{forumUsecase: forumUsecase_}
}

func (forumHandler *ForumHandler) CreateForum(c *gin.Context) {
	var forum models.Forum
	err := easyjson.UnmarshalFromReader(c.Request.Body, &forum)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
		return
	}

	err = forumHandler.forumUsecase.CreateForum(&forum)
	if err != nil && pkg.ConvertErrorToCode(err) != http.StatusConflict {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	if err != nil && pkg.ConvertErrorToCode(err) == http.StatusConflict {
		forumJson, internalErr := forum.MarshalJSON()
		if internalErr != nil {
			c.Data(pkg.CreateErrorResponse(err))
			return
		}
		c.Data(pkg.ConvertErrorToCode(err), "application/json; charset=utf-8", forumJson)
		return
	}

	forumJSON, err := forum.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", forumJSON)
}

func (forumHandler *ForumHandler) GetForum(c *gin.Context) {
	slug := c.Param("slug")

	forum, err := forumHandler.forumUsecase.GetInfoAboutForum(slug)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	forumJSON, err := forum.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", forumJSON)
}

func (forumHandler *ForumHandler) CreateThread(c *gin.Context) {
	slug := c.Param("slug")

	var thread models.Thread
	err := easyjson.UnmarshalFromReader(c.Request.Body, &thread)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
		return
	}
	thread.Forum = slug

	err = forumHandler.forumUsecase.CreateForumsThread(&thread)
	if err != nil && pkg.ConvertErrorToCode(err) != http.StatusConflict {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	if err != nil && pkg.ConvertErrorToCode(err) == http.StatusConflict {
		threadJson, internalErr := thread.MarshalJSON()
		if internalErr != nil {
			c.Data(pkg.CreateErrorResponse(err))
			return
		}
		c.Data(pkg.ConvertErrorToCode(err), "application/json; charset=utf-8", threadJson)
		return
	}
	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}
	c.Data(http.StatusCreated, "application/json; charset=utf-8", threadJSON)
}

func (forumHandler *ForumHandler) GetForumUsers(c *gin.Context) {
	slug := c.Param("slug")

	since := c.Query("since")

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

	users, err := forumHandler.forumUsecase.GetForumUsers(slug, defaultLimit, since, defaultDesc)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	usersJSON, err := users.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", usersJSON)
}

func (forumHandler *ForumHandler) GetForumThreads(c *gin.Context) {
	slug := c.Param("slug")

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
			return
		}
	}
	since := c.Query("since")
	descStr := c.Query("desc")
	desc := false
	if descStr != "" {
		var err error
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			c.Data(pkg.CreateErrorResponse(pkg.ErrBadRequest))
			return
		}
	}

	threads, err := forumHandler.forumUsecase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	threadsJSON, err := threads.MarshalJSON()
	if err != nil {
		c.Data(pkg.CreateErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadsJSON)
}
