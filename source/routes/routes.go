package routes

import (
	"chat/source/services"

	"github.com/gin-gonic/gin"
)
func Router (router *gin.Engine) {
	router.Static("/public", "./public")
    router.LoadHTMLGlob("templates/*")
    services.HandlerSSE(router)
}
