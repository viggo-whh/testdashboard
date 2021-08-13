package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK,"pong")
}

func WriteError(c *gin.Context, msg string)  {
	c.JSON(http.StatusOK,gin.H{
		"code":1,
		"msg":msg,
	})
}

func WriteOK(c *gin.Context, data interface{}) {
	ret, ok := data.(gin.H)
	if !ok {
		ret = gin.H{}
		ret["data"] = data
	}
	ret["code"] = 0
	ret["message"] = "success"
	c.JSON(http.StatusOK,ret)
}