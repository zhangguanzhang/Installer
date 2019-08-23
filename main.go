package main

import (
	"Installer/api"
	"Installer/api/v1"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//设定请求url不存在的返回值
	router.NoRoute(api.NoResponse)

	router.LoadHTMLGlob("templates/*")

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/ks", v1.GetKsFile)
		apiV1.POST("/ks", v1.KsUpdate)
		apiV1.POST("/upload", v1.ExcelUpload)
	}

	_ = router.Run(":8080")
}
