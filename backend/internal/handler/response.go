package handler

import (
	"errors"
	"github.com/blueship581/gbemr/internal/constants"
	"github.com/blueship581/gbemr/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{Code: constants.CodeOK, Message: "ok", Data: data})
}
func Fail(c *gin.Context, err error) {
	status, code := http.StatusInternalServerError, constants.CodeInternal
	msg := err.Error()
	switch {
	case errors.Is(err, repository.ErrNotFound):
		status, code = http.StatusNotFound, constants.CodeNotFound
	case msg == "invalid credentials":
		status, code = http.StatusUnauthorized, constants.CodeUnauthorized
	case msg == "forbidden":
		status, code = http.StatusForbidden, constants.CodeForbidden
	case msg == "username already exists":
		status, code = http.StatusConflict, constants.CodeConflict
	default:
		status, code = http.StatusBadRequest, constants.CodeBadRequest
	}
	c.JSON(status, response{Code: code, Message: msg})
}
