package controllers

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDataRecords 获取指定条记录，支持从指定条目开始计数
// @Summary 获取数据记录
// @Description 返回指定数量的记录，支持从指定偏移量开始
// @Tags data
// @Accept json
// @Produce json
// @Param limit query int false "返回的记录数量，默认为10"
// @Param offset query int false "从第几条记录开始，默认为0"
// @Success 200 {object} utils.JSONResponse{data=[]models.Data}
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data [get]
func GetDataRecords(c *gin.Context) {
	// 获取查询参数
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// 转换参数为整数
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		utils.JsonErrorResponse(c, 200400, "参数limit必须是正整数")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.JsonErrorResponse(c, 200400, "参数offset必须是非负整数")
		return
	}

	// 限制最大查询数量
	if limit > 100 {
		limit = 100
	}

	var data []models.Data
	db := database.GetDB()

	// 查询数据
	result := db.Offset(offset).Limit(limit).Find(&data)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 200500, "查询数据失败")
		return
	}

	// 获取总记录数
	var totalCount int64
	db.Model(&models.Data{}).Count(&totalCount)

	response := map[string]interface{}{
		"data":        data,
		"total":       totalCount,
		"limit":       limit,
		"offset":      offset,
		"has_more":    offset+limit < int(totalCount),
		"next_offset": offset + limit,
	}

	utils.JsonSuccessResponse(c, response)
}

// GetDataByID 根据ID获取单条数据记录
// @Summary 根据ID获取数据
// @Description 根据数据ID返回单条记录
// @Tags data
// @Accept json
// @Produce json
// @Param id path string true "数据ID"
// @Success 200 {object} utils.JSONResponse{data=models.Data}
// @Failure 400 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/{id} [get]
func GetDataByID(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.JsonErrorResponse(c, 200400, "参数id不能为空")
		return
	}

	var data models.PublicData
	db := database.GetDB()

	// 查询数据
	result := db.Table("data").Where("ID = ?", id).First(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			utils.JsonErrorResponse(c, 200404, "数据不存在")
		} else {
			utils.JsonErrorResponse(c, 200500, "查询数据失败")
		}
		return
	}

	utils.JsonSuccessResponse(c, data)
}

// GetDataStatistics 获取数据统计信息
// @Summary 获取数据统计
// @Description 返回化合物数量、物种数量、活性数据量、质谱数据量、核磁数据量等统计信息
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} utils.JSONResponse{data=map[string]interface{}}
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/statistics [get]
func GetDataStatistics(c *gin.Context) {
	db := database.GetDB()

	// 定义统计结果结构
	var stats struct {
		TotalCompounds  int64 `json:"total_compounds"`
		TotalSpecies    int64 `json:"total_species"`
		BioactivityData int64 `json:"bioactivity_data"`
		MSData          int64 `json:"ms_data"`
		NMRData         int64 `json:"nmr_data"`
	}

	// 获取总化合物数量（总记录数）
	db.Model(&models.Data{}).Count(&stats.TotalCompounds)

	// // 获取物种数量（固定值，这里假设是ItemType不为空的记录数）
	// db.Model(&models.Data{}).Where("ItemType IS NOT NULL AND ItemType != ''").Count(&stats.TotalSpecies)
	stats.TotalSpecies = 500

	// 获取活性数据量（Bioactivity不为空的数量）
	db.Model(&models.Data{}).Where("Bioactivity IS NOT NULL AND Bioactivity != ''").Count(&stats.BioactivityData)

	// 获取质谱数据量（MS1或MS2不为空的数量）
	db.Model(&models.Data{}).Where("MS1 IS NOT NULL OR MS2 IS NOT NULL AND MS2 != ''").Count(&stats.MSData)

	// 获取核磁数据量（NMR_13C_data不为空的数量）
	db.Model(&models.Data{}).Where("NMR_13C_data IS NOT NULL AND NMR_13C_data != ''").Count(&stats.NMRData)

	utils.JsonSuccessResponse(c, stats)
}

// FilterCompounds 筛选化合物
// @Summary 筛选化合物
// @Description 根据ItemType、分子量范围、Description和Source进行筛选，支持数组参数
// @Tags data
// @Accept json
// @Produce json
// @Param limit query int false "返回的记录数量，默认为10"
// @Param offset query int false "从第几条记录开始，默认为0"
// @Param item_type query []string false "ItemType分类数组" collectionFormat(multi)
// @Param min_weight query number false "最小分子量"
// @Param max_weight query number false "最大分子量"
// @Param description query []string false "Description描述数组" collectionFormat(multi)
// @Param source query []string false "Source来源数组" collectionFormat(multi)
// @Success 200 {object} utils.JSONResponse{data=[]models.Data}
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/filter [get]
func FilterCompounds(c *gin.Context) {
	// 获取查询参数
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	itemTypes := c.QueryArray("item_type")
	minWeightStr := c.Query("min_weight")
	maxWeightStr := c.Query("max_weight")
	descriptions := c.QueryArray("description")
	sources := c.QueryArray("source")

	// 转换参数为整数
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		utils.JsonErrorResponse(c, 200400, "参数limit必须是正整数")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.JsonErrorResponse(c, 200400, "参数offset必须是非负整数")
		return
	}

	// 限制最大查询数量
	if limit > 100 {
		limit = 100
	}

	// 转换分子量参数
	var minWeight, maxWeight float64

	if minWeightStr != "" {
		minWeight, err = strconv.ParseFloat(minWeightStr, 64)
		if err != nil {
			utils.JsonErrorResponse(c, 200400, "参数min_weight必须是数字")
			return
		}
	}

	if maxWeightStr != "" {
		maxWeight, err = strconv.ParseFloat(maxWeightStr, 64)
		if err != nil {
			utils.JsonErrorResponse(c, 200400, "参数max_weight必须是数字")
			return
		}
	}

	// 调用筛选服务，传入分页参数和数组参数
	compounds, totalCount, err := services.FilterCompounds(itemTypes, minWeight, maxWeight, descriptions, sources, limit, offset)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, "筛选化合物失败")
		return
	}

	response := map[string]interface{}{
		"data":        compounds,
		"total":       totalCount,
		"limit":       limit,
		"offset":      offset,
		"has_more":    offset+limit < int(totalCount),
		"next_offset": offset + limit,
	}

	utils.JsonSuccessResponse(c, response)
}

// GetItemTypes 获取所有ItemType分类
// @Summary 获取ItemType分类
// @Description 返回所有可用的ItemType分类
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} utils.JSONResponse{data=[]string}
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/item-types [get]
func GetItemTypes(c *gin.Context) {
	itemTypes, err := services.GetItemTypes()
	if err != nil {
		utils.JsonErrorResponse(c, 200500, "获取ItemType分类失败")
		return
	}

	utils.JsonSuccessResponse(c, itemTypes)
}

// GetDescriptions 获取所有Description分类
// @Summary 获取Description分类
// @Description 返回所有可用的Description分类
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} utils.JSONResponse{data=[]string}
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/descriptions [get]
func GetDescriptions(c *gin.Context) {
	descriptions, err := services.GetDescriptions()
	if err != nil {
		utils.JsonErrorResponse(c, 200500, "获取Description分类失败")
		return
	}

	utils.JsonSuccessResponse(c, descriptions)
}

// GetDataByIDFull 根据ID获取单条数据记录（保护数据）
// @Summary 根据ID获取数据（保护数据）
// @Description 根据数据ID返回MS2、Bioactivity和NMR_13C_data等保护数据
// @Tags data
// @Accept json
// @Produce json
// @Param id path string true "数据ID"
// @Success 200 {object} utils.JSONResponse{data=map[string]interface{}}
// @Failure 400 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/{id}/full [get]
func GetDataByIDFull(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.JsonErrorResponse(c, 200400, "参数id不能为空")
		return
	}

	var data models.ProtectedData
	db := database.GetDB()

	// 查询数据，只选择需要的字段
	result := db.Table("data").
		Select("MS2, Bioactivity, NMR_13C_data").
		Where("ID = ?", id).
		First(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			utils.JsonErrorResponse(c, 200404, "数据不存在")
		} else {
			utils.JsonErrorResponse(c, 200500, "查询数据失败")
		}
		return
	}

	utils.JsonSuccessResponse(c, data)
}

func GetMS2FullByID(c *gin.Context) {
	// 获取路径参数
	id := c.Param("id")
	if id == "" {
		utils.JsonErrorResponse(c, 200400, "参数id不能为空")
		return
	}

	var data string
	db := database.GetDB()

	// 查询数据，只选择需要的字段
	result := db.Table("data").
		Select("MS2_full").
		Where("ID = ?", id).
		First(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			utils.JsonErrorResponse(c, 200404, "数据不存在")
		} else {
			utils.JsonErrorResponse(c, 200500, "查询数据失败")
		}
		return
	}

	utils.JsonSuccessResponse(c, data)
}

// GetSources 获取所有Source分类
// @Summary 获取Source分类
// @Description 返回所有可用的Source分类
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} utils.JSONResponse{data=[]string}
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/sources [get]
func GetSources(c *gin.Context) {
	sources, err := services.GetSources()
	if err != nil {
		utils.JsonErrorResponse(c, 200500, "获取Source分类失败")
		return
	}

	utils.JsonSuccessResponse(c, sources)
}

// GetStructure 获取指定id的structure数据
// @Summary 获取指定id的structure数据
// @Description 获取指定id的structure数据
// @Tags data
// @Accept json
// @Produce json
// @Param id string
// @Success 200 {object} utils.JSONResponse{data=string}
// @Failure 500 {object} utils.JSONResponse
// @Router /api/data/sources [get]
func GetStructure(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.JsonErrorResponse(c, 200400, "参数id不能为空")
		return
	}
	// 定义只包含Structure的结构

	var data string
	db := database.GetDB()

	// 查询数据，只选择需要的字段
	result := db.Table("data").
		Select("Structure").
		Where("ID = ?", id).
		First(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			utils.JsonErrorResponse(c, 200404, "数据不存在")
		} else {
			utils.JsonErrorResponse(c, 200500, "查询数据失败")
		}
		return
	}

	utils.JsonSuccessResponse(c, data)
}
