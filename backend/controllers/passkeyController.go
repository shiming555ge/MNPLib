package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PasskeyRequest passkey 请求结构
type PasskeyRequest struct {
	Description string `json:"description"`
	Operator    string `json:"operator" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

// PasskeyResponse passkey 响应结构
type PasskeyResponse struct {
	Passkey     string    `json:"passkey"`
	Description string    `json:"description"`
	Operator    string    `json:"operator"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
	Extends     string    `json:"extends"`
}

// GetAllPasskeys 获取所有 passkey
// @Summary 获取所有 passkey
// @Description 获取所有 passkey 列表，需要管理员权限。注意：不会返回当前登录用户自己的 passkey
// @Tags passkey
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.JSONResponse{data=[]PasskeyResponse}
// @Failure 401 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/passkeys [get]
func GetAllPasskeys(c *gin.Context) {
	// 获取当前登录用户的 passkey
	currentPasskey, exists := c.Get("passkey")
	if !exists {
		utils.JsonErrorResponse(c, 200401, "未获取到用户信息")
		return
	}
	currentPasskeyStr, _ := currentPasskey.(string)

	var passkeys []models.Passkey
	db := database.GetDB()

	// 查询所有 passkey，排除当前用户的 passkey
	result := db.Where("Passkey != ?", currentPasskeyStr).Order("Created_At DESC").Find(&passkeys)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 200500, "数据库查询失败")
		return
	}

	// 转换为响应格式
	var response []PasskeyResponse
	for _, p := range passkeys {
		response = append(response, PasskeyResponse{
			Passkey:     p.Passkey,
			Description: p.Description,
			Operator:    p.Operator,
			CreatedAt:   p.CreatedAt,
			IsActive:    p.IsActive,
			Extends:     p.Extends,
		})
	}

	utils.JsonSuccessResponse(c, response)
}

// CreatePasskey 创建新的 passkey
// @Summary 创建 passkey
// @Description 创建新的 passkey，系统会自动生成 UUID。创建者的 operator(passkey) 信息会被存储在 Extends 字段中
// @Tags passkey
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body PasskeyRequest true "passkey 信息"
// @Success 201 {object} utils.JSONResponse{data=PasskeyResponse}
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/passkeys [post]
func CreatePasskey(c *gin.Context) {
	var req PasskeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 200400, "请求参数错误")
		return
	}

	// 从上下文中获取创建者的 operator 和 passkey 信息
	creatorOperator, exists := c.Get("operator")
	if !exists {
		creatorOperator = "未知"
	}
	creatorOperatorStr, _ := creatorOperator.(string)

	creatorPasskey, exists := c.Get("passkey")
	if !exists {
		creatorPasskey = "未知"
	}
	creatorPasskeyStr, _ := creatorPasskey.(string)

	// 构建 Extends 信息，格式为 "operator(passkey)"
	extendsInfo := creatorOperatorStr + "(" + creatorPasskeyStr + ")"

	// 创建 passkey（数据库会自动生成 UUID）
	passkey := models.Passkey{
		Description: req.Description, // description 可以为空
		Operator:    req.Operator,
		IsActive:    req.IsActive,
		Extends:     extendsInfo,
	}

	db := database.GetDB()
	result := db.Create(&passkey)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 200500, "创建 passkey 失败")
		return
	}

	response := PasskeyResponse{
		Passkey:     passkey.Passkey,
		Description: passkey.Description,
		Operator:    passkey.Operator,
		CreatedAt:   passkey.CreatedAt,
		IsActive:    passkey.IsActive,
		Extends:     passkey.Extends,
	}

	utils.JsonSuccessResponse(c, response)
}

// UpdatePasskey 更新 passkey
// @Summary 更新 passkey
// @Description 更新指定 passkey 的信息。注意：不能修改自己的 passkey
// @Tags passkey
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param passkey path string true "passkey UUID"
// @Param request body PasskeyRequest true "更新信息"
// @Success 200 {object} utils.JSONResponse{data=PasskeyResponse}
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/passkeys/{passkey} [put]
func UpdatePasskey(c *gin.Context) {
	passkeyID := c.Param("passkey")
	if passkeyID == "" {
		utils.JsonErrorResponse(c, 200400, "passkey 参数不能为空")
		return
	}

	// 获取当前登录用户的 passkey
	currentPasskey, exists := c.Get("passkey")
	if !exists {
		utils.JsonErrorResponse(c, 200401, "未获取到用户信息")
		return
	}
	currentPasskeyStr, _ := currentPasskey.(string)

	// 检查是否尝试修改自己的 passkey
	if passkeyID == currentPasskeyStr {
		utils.JsonErrorResponse(c, 200403, "不能修改自己的 passkey")
		return
	}

	var req PasskeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 200400, "请求参数错误")
		return
	}

	db := database.GetDB()
	var passkey models.Passkey

	// 查找 passkey
	result := db.Where("Passkey = ?", passkeyID).First(&passkey)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200404, "passkey 不存在")
		} else {
			utils.JsonErrorResponse(c, 200500, "数据库查询失败")
		}
		return
	}

	// 更新信息
	passkey.Description = req.Description
	passkey.Operator = req.Operator
	passkey.IsActive = req.IsActive

	result = db.Save(&passkey)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 200500, "更新 passkey 失败")
		return
	}

	response := PasskeyResponse{
		Passkey:     passkey.Passkey,
		Description: passkey.Description,
		Operator:    passkey.Operator,
		CreatedAt:   passkey.CreatedAt,
		IsActive:    passkey.IsActive,
		Extends:     passkey.Extends,
	}

	utils.JsonSuccessResponse(c, response)
}

// DeletePasskey 删除 passkey
// @Summary 删除 passkey
// @Description 删除指定的 passkey。注意：不能删除自己的 passkey
// @Tags passkey
// @Security BearerAuth
// @Produce json
// @Param passkey path string true "passkey UUID"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/passkeys/{passkey} [delete]
func DeletePasskey(c *gin.Context) {
	passkeyID := c.Param("passkey")
	if passkeyID == "" {
		utils.JsonErrorResponse(c, 200400, "passkey 参数不能为空")
		return
	}

	// 获取当前登录用户的 passkey
	currentPasskey, exists := c.Get("passkey")
	if !exists {
		utils.JsonErrorResponse(c, 200401, "未获取到用户信息")
		return
	}
	currentPasskeyStr, _ := currentPasskey.(string)

	// 检查是否尝试删除自己的 passkey
	if passkeyID == currentPasskeyStr {
		utils.JsonErrorResponse(c, 200403, "不能删除自己的 passkey")
		return
	}

	db := database.GetDB()
	result := db.Where("Passkey = ?", passkeyID).Delete(&models.Passkey{})
	if result.Error != nil {
		utils.JsonErrorResponse(c, 200500, "删除 passkey 失败")
		return
	}

	if result.RowsAffected == 0 {
		utils.JsonErrorResponse(c, 200404, "passkey 不存在")
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// GetPasskeyByID 获取单个 passkey 信息
// @Summary 获取单个 passkey
// @Description 根据 passkey UUID 获取单个 passkey 的详细信息
// @Tags passkey
// @Security BearerAuth
// @Produce json
// @Param passkey path string true "passkey UUID"
// @Success 200 {object} utils.JSONResponse{data=PasskeyResponse}
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/passkeys/{passkey} [get]
func GetPasskeyByID(c *gin.Context) {
	passkeyID := c.Param("passkey")
	if passkeyID == "" {
		utils.JsonErrorResponse(c, 200400, "passkey 参数不能为空")
		return
	}

	db := database.GetDB()
	var passkey models.Passkey

	// 查找 passkey
	result := db.Where("Passkey = ?", passkeyID).First(&passkey)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200404, "passkey 不存在")
		} else {
			utils.JsonErrorResponse(c, 200500, "数据库查询失败")
		}
		return
	}

	response := PasskeyResponse{
		Passkey:     passkey.Passkey,
		Description: passkey.Description,
		Operator:    passkey.Operator,
		CreatedAt:   passkey.CreatedAt,
		IsActive:    passkey.IsActive,
		Extends:     passkey.Extends,
	}

	utils.JsonSuccessResponse(c, response)
}

// TogglePasskeyStatus 切换 passkey 状态
// @Summary 切换 passkey 状态
// @Description 启用或禁用指定的 passkey。注意：不能修改自己的 passkey 状态
// @Tags passkey
// @Security BearerAuth
// @Produce json
// @Param passkey path string true "passkey UUID"
// @Success 200 {object} utils.JSONResponse{data=PasskeyResponse}
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/passkeys/{passkey}/toggle [post]
func TogglePasskeyStatus(c *gin.Context) {
	passkeyID := c.Param("passkey")
	if passkeyID == "" {
		utils.JsonErrorResponse(c, 200400, "passkey 参数不能为空")
		return
	}

	// 获取当前登录用户的 passkey
	currentPasskey, exists := c.Get("passkey")
	if !exists {
		utils.JsonErrorResponse(c, 200401, "未获取到用户信息")
		return
	}
	currentPasskeyStr, _ := currentPasskey.(string)

	// 检查是否尝试修改自己的 passkey
	if passkeyID == currentPasskeyStr {
		utils.JsonErrorResponse(c, 200403, "不能修改自己的 passkey 状态")
		return
	}

	db := database.GetDB()
	var passkey models.Passkey

	// 查找 passkey
	result := db.Where("Passkey = ?", passkeyID).First(&passkey)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200404, "passkey 不存在")
		} else {
			utils.JsonErrorResponse(c, 200500, "数据库查询失败")
		}
		return
	}

	// 切换状态
	passkey.IsActive = !passkey.IsActive

	result = db.Save(&passkey)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 200500, "更新 passkey 状态失败")
		return
	}

	response := PasskeyResponse{
		Passkey:     passkey.Passkey,
		Description: passkey.Description,
		Operator:    passkey.Operator,
		CreatedAt:   passkey.CreatedAt,
		IsActive:    passkey.IsActive,
		Extends:     passkey.Extends,
	}

	utils.JsonSuccessResponse(c, response)
}
