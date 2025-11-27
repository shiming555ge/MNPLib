package controllers

import (
	"backend/database"
	"backend/models"
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

	var data models.Data
	db := database.GetDB()

	// 查询数据
	result := db.Where("ID = ?", id).First(&data)
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
