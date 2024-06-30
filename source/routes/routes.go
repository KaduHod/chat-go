package routes
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router (router *gin.Engine) {
	router.GET("/", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message": "success",
		})
	})
	router.GET("/home", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message": "success",
		})		
	})
}
