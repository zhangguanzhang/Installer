package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	NotFound = errors.New("NotFound")
)



// Response setting gin.JSON
func NewResponse(ctx *gin.Context, httpCode int, Msg, data interface{}) {
	var msg  = "ok"

	switch v := Msg.(type) {
	case error:
		msg = v.Error()
	case string:
		msg = v
	}

	ctx.JSON(httpCode, Response{
		Code:    httpCode,
		Message: msg,
		Data: data,
	})
}


type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}


//返回404
func NoResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": http.StatusNotFound,
		"error":  "404, page not exists!",
	})
}

