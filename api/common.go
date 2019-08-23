package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoResponse(c *gin.Context) {
	//返回404状态码
	c.JSON(http.StatusNotFound, gin.H{
		"status": http.StatusNotFound,
		"error":  "404, page not exists!",
	})
}
