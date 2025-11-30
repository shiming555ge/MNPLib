package router

import (
	"backend/config"
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
			data.GET("/statistics", controllers.GetDataStatistics)
			data.GET("/filter", controllers.FilterCompounds)
			data.GET("/item-types", controllers.GetItemTypes)
			data.GET("/descriptions", controllers.GetDescriptions)
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
		static := r.Group("/")
		{
			static.Static("/", "./www")
		}
	}
}
