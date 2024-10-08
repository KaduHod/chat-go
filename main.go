package main

import (
	"chat/source/routes"
	"chat/source/utils"
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
)
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        // Stop timer
        end := time.Now()
        latency := end.Sub(start)

        // Log request details
        status := c.Writer.Status()
        clientIP := c.ClientIP()
        method := c.Request.Method
        path := c.Request.URL.Path

        log := fmt.Sprintf("%v | %3d | %8v | %8s | %-5s %s\n",
        end.Format(utils.DDMMYYY_hhmmss),
        status,
        latency,
        clientIP,
        method,
        path,
    )

    fmt.Printf(log)
    utils.Logger("/request.log", log, "REQUISICAO", false)
    c.Next()
    }
}
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
func main() {
	app := gin.Default()
	app.Use(Logger())
    app.Use(CORSMiddleware())
	app.Use(gin.Recovery())
	routes.Router(app)
	app.Run(":3000")
}
