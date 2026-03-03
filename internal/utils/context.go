package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

// ---------------------------
// @brief 从上下文中获取用户ID
// ---------------------------
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userIDVal, exists := c.Get(UserIDKey)
	if !exists {
		return 0, errors.New("上下文中未找到用户ID")
	}

	// 类型断言
	userID, ok := userIDVal.(uint)
	if !ok {
		return 0, errors.New("用户ID类型错误")
	}

	return userID, nil
}
