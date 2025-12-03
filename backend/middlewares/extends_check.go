package middlewares

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExtendsCheck 检查用户 Extends 字段是否为空的中间件
// 只有 Extends 为空的用户才能访问某些功能
func ExtendsCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取当前用户的 passkey
		currentPasskey, exists := c.Get("passkey")
		if !exists {
			utils.JsonErrorResponse(c, http.StatusUnauthorized, "未获取到用户信息")
			c.Abort()
			return
		}
		currentPasskeyStr, _ := currentPasskey.(string)

		// 查询数据库获取用户的 Extends 信息
		db := database.GetDB()
		var passkey models.Passkey
		result := db.Where("Passkey = ?", currentPasskeyStr).First(&passkey)
		if result.Error != nil {
			utils.JsonErrorResponse(c, http.StatusInternalServerError, "查询用户信息失败")
			c.Abort()
			return
		}

		// 检查 Extends 字段是否为空
		if passkey.Extends != "" {
			utils.JsonErrorResponse(c, http.StatusForbidden, "只有超级管理员才能访问此功能")
			c.Abort()
			return
		}

		c.Next()
	}
}
