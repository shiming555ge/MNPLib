package middlewares

import (
	"backend/config"
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.JsonErrorResponse(c, http.StatusUnauthorized, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JsonErrorResponse(c, http.StatusUnauthorized, "令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析和验证token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Config.GetString("jwt.secret")), nil
		})

		if err != nil {
			utils.JsonErrorResponse(c, http.StatusUnauthorized, "令牌无效或已过期")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 获取用户信息
			passkey, _ := claims["passkey"].(string)
			operator, _ := claims["operator"].(string)
			description, _ := claims["description"].(string)

			// 记录访问日志
			utils.LogAccess(passkey, operator, description, c.Request.Method, c.Request.URL.Path, c.ClientIP())

			// 将用户信息存储到上下文中
			c.Set("passkey", passkey)
			c.Set("operator", operator)
			c.Set("description", description)
			c.Next()
		} else {
			utils.JsonErrorResponse(c, http.StatusUnauthorized, "令牌无效")
			c.Abort()
			return
		}
	}
}
