package middleware

// JWTAuth 鉴权中间件，AdminAuth 管理员权限中间件，保护需要登录和管理员权限的接口

import (
	"net/http"

	"geekedu-project/common/jwt"

	"github.com/gin-gonic/gin"
)

// JWTAuth 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取Token（从Header的Authorization字段）
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			c.Abort()
			return
		}

		// 关键修复：剔除"Bearer "前缀（前端发送 "Bearer xxxxx" 格式）
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 2. 解析Token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token无效",
			})
			c.Abort()
			return
		}

		// 3. 将用户信息存入Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminAuth 管理员权限中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "无管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
