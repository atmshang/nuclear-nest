package versionutil

import (
	"github.com/atmshang/nuclear-nest/pkg/apiutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetVersionInfoFunc(c *gin.Context) {
	c.JSON(http.StatusOK, apiutil.Response{
		Code:    2000,
		Message: "",
		Data:    GetVersionInfo,
	})
}
