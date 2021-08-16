package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testdash/controllers"
)

func InitApi(eng *gin.Engine) {
	// gin配置使用中间件
	//eng.Use(CorsMiddleware)

	//health check 接口
	eng.GET("/ping", controllers.Ping)

	//接口分组
	api := eng.Group("/api/v1")
	api.GET("nodes", controllers.GetNodeList)

	//获取metric指标数据
	api.POST("metrics", controllers.GetMetrics)

	api.GET("namespaces/:namespace/pods/:pod/logs", controllers.GetKubeLogs)
}

func CorsMiddleware(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, UPDATA")
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	c.Header("add_header Access-Control-Allow-Credentials", "true")
	c.Set("content-type", "aoolication/json")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
