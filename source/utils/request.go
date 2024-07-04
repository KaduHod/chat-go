package utils
import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
func RequestBody[T any](bodyDest * T, c *gin.Context) bool  {
	if err := c.BindJSON(bodyDest); err != nil {
		log.Println("Erro :: pegando body do request")
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H {
			"mensagem": "corpo de requisição é necessário",
		})
		return false
	}
	return true
}
