package security

import (
    "strings"
    
    "github.com/gin-gonic/gin"
)

type CORSConfig struct {
    AllowOrigins []string
    AllowMethods []string
    AllowHeaders []string
}

func NewCORSMiddleware(config CORSConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 设置CORS头
        c.Header("Access-Control-Allow-Origin", strings.Join(config.AllowOrigins, ","))
        c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
        c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
        c.Header("Access-Control-Allow-Credentials", "true")

        // 处理预检请求
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }

        c.Next()
    }
}