package controllers

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Passkey string `json:"passkey" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Operator    string `json:"operator"`
	Token       string `json:"token"`
	Description string `json:"description"`
}

// Login 登录API
// @Summary 用户登录
// @Description 使用passkey进行登录，返回JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} utils.JSONResponse{data=LoginResponse}
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	// 解析请求体
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 200400, "请求参数错误")
		return
	}

	// 验证passkey是否为空
	if req.Passkey == "" {
		utils.JsonErrorResponse(c, 200400, "passkey不能为空")
		return
	}

	// 查询数据库验证passkey
	var passkey models.Passkey
	db := database.GetDB()

	result := db.Where("Passkey = ? AND Is_Active = ?", req.Passkey, true).First(&passkey)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200401, "无效的passkey或已禁用")
		} else {
			utils.JsonErrorResponse(c, 200500, "数据库查询失败")
		}
		return
	}

	// 创建JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passkey":     req.Passkey,
		"operator":    passkey.Operator,
		"description": passkey.Description,
		"exp":         time.Now().Add(time.Hour * 24 * 3).Unix(), // 3天过期
		"iat":         time.Now().Unix(),
	})

	// 获取JWT secret
	secret := config.Config.GetString("jwt.secret")
	if secret == "" {
		utils.JsonErrorResponse(c, 200500, "服务器配置错误")
		return
	}

	// 签名token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		utils.JsonErrorResponse(c, 200500, "生成token失败")
		return
	}

	// 返回token
	response := LoginResponse{
		Token:       tokenString,
		Operator:    passkey.Operator,
		Description: passkey.Description,
	}

	utils.JsonSuccessResponse(c, response)
}

// VerifyLoginStatus 验证登录状态
// @Summary 验证登录状态
// @Description 验证JWT token是否有效（中间件已验证），返回成功状态
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Router /api/auth/verify [get]
func VerifyLoginStatus(c *gin.Context) {
	// 中间件已经验证了JWT token，这里只需要返回成功状态
	utils.JsonSuccessResponse(c, nil)
}

// VerifyPasskeyModifiable 验证是否可以修改passkey
// @Summary 验证是否可以修改passkey
// @Description 验证当前用户是否有权限修改passkey（中间件已验证Extends字段为空），返回成功状态
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Router /api/auth/verify-passkey-modifiable [get]
func VerifyPasskeyModifiable(c *gin.Context) {
	// 中间件已经验证了JWT token和Extends字段为空，这里只需要返回成功状态
	utils.JsonSuccessResponse(c, nil)
}
