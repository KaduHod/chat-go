package utils

import (
	"bytes"
	"io"
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
func RestPostRequest(baseUrl string, dados []byte) ([]byte, *http.Response, error) {
    respostaHttp, err := http.Post(baseUrl, "application/json", bytes.NewBuffer(dados))
    if err != nil {
        return nil, nil, err
    }
    defer respostaHttp.Body.Close()
    corpoRespostaHttp, err := io.ReadAll(respostaHttp.Body)
    if err != nil {
        return nil, nil, err
    }
    return corpoRespostaHttp, respostaHttp, nil
}
func RestGetRequest(baseUrl string) ([]byte, *http.Response, error) {
    respostaHttp, err := http.Get(baseUrl)
    if err != nil {
        return nil, nil, err
    }
    defer respostaHttp.Body.Close()
    corpoRespostaHttp, err := io.ReadAll(respostaHttp.Body)
    if err != nil {
        return nil, nil, err
    }
    return corpoRespostaHttp, respostaHttp, nil
}
