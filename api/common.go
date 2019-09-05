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
func NewResponse(ctx *gin.Context, httpCode int, info ...interface{}) {
	var (
		msg  = "ok"
		data interface{} = nil
	)

	if httpCode != http.StatusOK {
		msg = "not ok"
	}
	if len(info) >= 1 && info[0] != nil {
		switch v := info[0].(type) {
		case error:
			msg = v.Error()
		case string:
			msg = v
		}
	}
	if len(info) >= 2 && info[1] != nil {
		data = info[1]
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

