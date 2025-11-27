package router

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// API路由组
	api := r.Group("/api")
	{
		// 数据相关路由
		data := api.Group("/data")
		{
			data.GET("", controllers.GetDataRecords)
			data.GET("/:id", controllers.GetDataByID)
		}
	}
}
