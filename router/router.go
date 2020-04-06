package router

import (
	"Installer/api"
	v1 "Installer/api/v1"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func InitRouter(KSTemplate string) *gin.Engine {

	router := gin.Default()

	//设定请求url不存在的返回值
	router.NoRoute(api.NoResponse)

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency, //时间
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	apiV1 := router.Group("/api/v1")
	{

		apiV1.GET("/ks", v1.GetKsFile(KSTemplate))
		apiV1.POST("/ks", v1.UpdateStatusFromKs)

		apiV1.GET("/status", v1.GetStatus)

		apiV1.GET("/machines", v1.GetMachines)

		apiV1.GET("/sns", v1.GetSNS)

		apiV1.GET("/machine/:sn", v1.GetMachine)
		apiV1.POST("/machine", v1.AddMachine)
		apiV1.PUT("/machine/:sn", v1.UpdateMachine)
		apiV1.DELETE("/machine/:sn", v1.DeleteMachine)

		apiV1.POST("/upload", v1.UploadExcel)
	}

	return router

}
