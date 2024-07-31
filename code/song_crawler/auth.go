package song_crawler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 中间件函数,验证 JWT 是否有效
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization 字段
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		// 解析和验证 JWT
		token, err := jwt.ParseWithClaims(authHeader, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return JWTKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid JWT token",
			})
			c.Abort()
			return
		}

		// 将解析后的 claims 信息保存在 gin.Context 中,供后续处理使用
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid JWT claims",
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)

		// 继续处理请求
		c.Next()
	}
}

// generateJWTToken 生成 JWT token
func GenerateJWTToken(payload map[string]interface{}, secretKey []byte, expirationTime time.Duration) (string, error) {
	// 创建 JWT claims
	claims := jwt.MapClaims{}
	for key, value := range payload {
		claims[key] = value
	}

	// 设置过期时间
	claims["exp"] = time.Now().Add(expirationTime).Unix()

	// 创建 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
