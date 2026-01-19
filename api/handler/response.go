package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fail(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func badRequest(c *gin.Context, msg string) {
	fail(c, http.StatusBadRequest, 400, msg)
}

func internalError(c *gin.Context, msg string) {
	fail(c, http.StatusInternalServerError, 500, msg)
}

