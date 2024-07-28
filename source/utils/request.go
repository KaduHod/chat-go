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
func RestGetRequestAutenticado(baseUrl string, authToken string) ([]byte, *http.Response, error) {
    // Cria uma nova requisição
    req, err := http.NewRequest("GET", baseUrl, nil)
    if err != nil {
        return nil, nil, err
    }

    // Adiciona o header de autenticação
    req.Header.Add("Authorization", "Bearer " + authToken)

    // Executa a requisição
    client := &http.Client{}
    respostaHttp, err := client.Do(req)
    if err != nil {
        return nil, nil, err
    }
    defer respostaHttp.Body.Close()

    // Lê o corpo da resposta
    corpoRespostaHttp, err := io.ReadAll(respostaHttp.Body)
    if err != nil {
        return nil, nil, err
    }

    return corpoRespostaHttp, respostaHttp, nil
}
