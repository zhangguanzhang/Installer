package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	NotFound = errors.New("NotFound")
)

type CommonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Response setting gin.JSON
func Response(ctx *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	ctx.JSON(httpStatus, CommonResp{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func Success(ctx *gin.Context, data interface{}, msg string) {
	Response(ctx, http.StatusOK, 200, data, msg)
}

func Error(ctx *gin.Context, httpStatus int, msg string) {
	Response(ctx, httpStatus, httpStatus, nil, msg)
}

//ks有错误的话必须返回非200状态码,请勿使用此回应ks
//状态码为ok，但是回应错误信息
func Fail(ctx *gin.Context, msg string) {
	Response(ctx, http.StatusOK, 200, nil, msg)
}

//返回404
func NoResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": http.StatusNotFound,
		"error":  "404, page not exists!",
	})
}
