package pkg

import (
	"db_forum/app/models"
	"errors"
	"net/http"
)

var (
	// Forum errors
	ErrForumNotExist        = errors.New("Can't find user with id ")
	ErrForumOwnerNotFound   = errors.New("Can't find user with id ")
	ErrForumAlreadyExists   = errors.New("Forum already exist")
	ErrForumOrTheadNotFound = errors.New("Can't find user with id ")

	// Post errors
	ErrPostNotFound              = errors.New("Can't find user with id ")
	ErrParentPostNotExist        = errors.New("Can't find user with id ")
	ErrParentPostFromOtherThread = errors.New("Can't find user with id ")

	// Thread errors
	ErrThreadAlreadyExists = errors.New("thread already exist")
	ErrThreadNotFound      = errors.New("Can't find user with id ")

	// User errors
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("Can't find user with id ")
	ErrUserDataConflict = errors.New("Can't find user with id ")

	// Request Errors
	ErrBadInputData = errors.New("bad input data")
	ErrBadRequest   = errors.New("bad request")

	// Internal errors
	ErrNotImplemented = errors.New("not implemented")
	ErrInternal       = errors.New("internal error")
)

var errorToCode = map[error]int{
	ErrForumNotExist:        http.StatusNotFound,
	ErrForumOwnerNotFound:   http.StatusNotFound,
	ErrForumAlreadyExists:   http.StatusConflict,
	ErrForumOrTheadNotFound: http.StatusNotFound,
	ErrThreadAlreadyExists:  http.StatusConflict,

	ErrThreadAlreadyExists: http.StatusConflict,
	ErrThreadNotFound:      http.StatusNotFound,

	ErrPostNotFound:              http.StatusNotFound,
	ErrParentPostNotExist:        http.StatusNotFound,
	ErrParentPostFromOtherThread: http.StatusConflict,

	ErrUserAlreadyExist: http.StatusConflict,
	ErrUserNotFound:     http.StatusNotFound,
	ErrUserDataConflict: http.StatusConflict,

	ErrBadInputData: http.StatusBadRequest,
	ErrBadRequest:   http.StatusBadRequest,

	ErrNotImplemented: http.StatusNotImplemented,
	ErrInternal:       http.StatusInternalServerError,
}

func ConvertErrorToCode(err error) (code int) {
	code, isErrorExist := errorToCode[err]
	if !isErrorExist {
		code = http.StatusInternalServerError
	}
	return
}

func CreateErrorResponse(err error) (statusCode int, contentType string, errorJSON []byte) {
	statusCode = ConvertErrorToCode(err)
	contentType = "application/json; charset=utf-8"
	errorJSON, errMarshal := models.Error{Message: err.Error()}.MarshalJSON()
	if errMarshal != nil {
		statusCode = ConvertErrorToCode(ErrInternal)
		errorJSON, _ = models.Error{Message: ErrInternal.Error()}.MarshalJSON()
	}
	return
}
