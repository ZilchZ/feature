package middlewares

import (
    "github.com/feature/sdk/log"
    "github.com/gin-gonic/gin"
    "time"
)

func LogAccess(start time.Time, c *gin.Context) {
    latency := time.Since(start).Seconds()
    path := c.Request.URL.Path
    raw := c.Request.URL.RawQuery
    clientIP := c.ClientIP()
    method := c.Request.Method
    statusCode := c.Writer.Status()
    if raw != "" {
        path += "?" + raw
    }
    if statusCode >= 400 {
        log.ApiLogger.Errorf("[%s] %3d | %.3f | %s | %s", method, statusCode, latency, clientIP, path)
    } else {
        log.ApiLogger.Infof("[%s] %3d | %.3f | %s | %s", method, statusCode, latency, clientIP, path)
    }
    if len(c.Errors) > 0 {
        for _, v := range c.Errors {
            log.Logger.Error(v)
        }
    }
    
}
