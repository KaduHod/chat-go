package main

import (
	"chat/source/routes"
	"chat/source/services"
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

func main() {
    services.VerificaCanaisAoIniciarServidor()
    services.VerificaClientesEmCanaisAoIniciarServidor()
	app := gin.Default()
	app.Use(Logger())
	app.Use(gin.Recovery())
	routes.Router(app)
	app.Run(":3000")
}
