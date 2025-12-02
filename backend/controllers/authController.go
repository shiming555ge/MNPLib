package controllers

import (
	"backend/config"
	"backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Passkey string `json:"passkey" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string `json:"token"`
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

	// 这里可以添加passkey的验证逻辑
	// 例如：检查passkey是否在数据库中存在，或者是否符合某种规则
	// 目前我们假设任何非空的passkey都是有效的

	// 创建JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passkey": req.Passkey,
		"exp":     time.Now().Add(time.Hour * 24 * 3).Unix(), // 3天过期
		"iat":     time.Now().Unix(),
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
		Token: tokenString,
	}

	utils.JsonSuccessResponse(c, response)
}
