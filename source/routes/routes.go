package routes

import (
	"chat/source/services"
	"net/http"

	"github.com/gin-gonic/gin"
)
func Router (router *gin.Engine) {
	services.IniciarCanalPadrao()
	router.GET("/", func (c *gin.Context) {
		c.Header("Content-type", "text/html")
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.POST("/logar", services.LogarUsuario)
	router.GET("/ws", services.WebsocketHandler)
}
