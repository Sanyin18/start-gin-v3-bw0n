package song_crawler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c *gin.Context) {
	// 解析请求body中的JSON数据
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证用户身份
	if err := validateUser(loginRequest.Username, loginRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	payload := map[string]interface{}{
		"username": loginRequest.Username,
	}
	// 生成JWT token
	token, err := GenerateJWTToken(payload, JWTKey, 1*time.Hour)
	fmt.Println(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 返回登录响应
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}

func validateUser(username, password string) error {
	// 在数据库中查找用户
	storedPassword, ok := AdminUser[username]
	if !ok || storedPassword != password {
		return fmt.Errorf("invalid username or password")
	}
	return nil
}
