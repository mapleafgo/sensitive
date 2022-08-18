package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CheckAndPrint 校验请求头
func CheckAndPrint() gin.HandlerFunc {
	return func(c *gin.Context) {
		ct := c.GetHeader("Content-Type")
		if strings.Contains(ct, "application/json") && (strings.Contains(ct, "charset=utf-8") || strings.Contains(ct, "charset=UTF-8")) {
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Content-Type must be application/json charset=utf-8",
			})
			c.Abort()
		}
	}
}
