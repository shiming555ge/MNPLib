package router

import (
	"backend/config"
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
		}

		// 数据相关路由
		data := api.Group("/data")
		{
			// data.GET("", controllers.GetDataRecords)
			data.GET("/:id", controllers.GetDataByID)
			data.GET("/:id/structure", controllers.GetStructure)
			data.GET("/statistics", controllers.GetDataStatistics)
			data.GET("/filter", controllers.FilterCompounds)
			data.GET("/item-types", controllers.GetItemTypes)
			data.GET("/descriptions", controllers.GetDescriptions)
			data.GET("/sources", controllers.GetSources)
			// 受保护的数据路由，需要JWT认证
			data.GET("/:id/full", middlewares.JWTAuth(), controllers.GetDataByIDFull)
		}
		// RDKit相关路由
		rdkit := api.Group("/rdkit")
		{
			rdkit.GET("/status", controllers.GetRdkitStatus)
			rdkit.GET("/similarity", controllers.SimilaritySearch)
			rdkit.GET("/smiles-to-fingerprint", controllers.SmilesToFingerprint)
			rdkit.GET("/smiles-to-pdb", controllers.SmilesToPDB)
			rdkit.GET("/is-substructure", controllers.IsSubstructure)
			rdkit.GET("/substructure-search", controllers.SubstructureSearch)
			rdkit.GET("/exact-match", controllers.ExactMatchSearch)
		}
	}

	if config.Config.GetBool("static") {
		// 静态资源
		r.Static("/assets", "./www/assets")
		r.Static("/home_pics", "./www/home_pics")
		r.Static("/Ketcher", "./www/Ketcher")
		r.StaticFile("/logo.png", "./www/logo.png")
		r.StaticFile("/vite.svg", "./www/vite.svg")
		r.StaticFile("/favicon.ico", "./www/favicon.ico")

		// 根路径
		r.GET("/", func(c *gin.Context) {
			c.File("./www/index.html")
		})

		// 处理 SPA 客户端路由
		r.NoRoute(func(c *gin.Context) {
			c.File("./www/index.html")
		})
	}
}
