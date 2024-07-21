package routes

import (
	"chat/source/modules/financas"
	"chat/source/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)
func Router (router *gin.Engine) {
	services.IniciarCanalPadrao()
	router.GET("/", func (c *gin.Context) {
		c.Header("Content-type", "text/html")
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/financas/grafico", func (c* gin.Context) {
		c.HTML(http.StatusOK, "grafico.html", gin.H{})
	})
	router.POST("/logar", services.LogarUsuario)
    router.GET("/chat/canais", services.ListarCanaisOnlineHandler)
    router.GET("/chat/criar", services.CriarCanalHandler)
    router.GET("/chat/canal/:id/iniciar", services.IniciarCanalHandler)
    router.GET("/chat/canal/:id/fechar", services.FecharCanalHandler)
    router.GET("/chat/usuario/:id/canal/:idcanal/entrar", services.AdicionarUsuarioCanalHandler)
    router.GET("/chat/usuario/:id/canal/:idcanal/sair", services.RemoverUsuarioCanalHandler)
    router.GET("/ws/canal/:idcanal/cliente/:idusuario", services.WebsocketHandlerV2)
	router.GET("/ws", services.WebsocketHandler)
	router.GET("/financas/juroscomposto/simular",func(c *gin.Context) {
		var aplicacao financas.AplicacaoFinanceira
		var aporte financas.AporteAplicacaoFinanceira
		var erros []string
		quantidadeDeMeses, err := strconv.Atoi(c.Query("quantidadeDeMeses"))
		if err != nil {
			erros = append(erros, "quantidadeDeMeses inválido")
		}

		valorInicial, err := strconv.ParseFloat(c.Query("valorInicial"), 64)
		if err != nil {
			erros = append(erros, "valorInicial inválido")
		}

		taxa, err := strconv.ParseFloat(c.Query("taxa"), 64)
		if err != nil {
			erros = append(erros, "taxa inválida")
		}

		valorAporte, err := strconv.ParseFloat(c.Query("valorAporte"), 64)
		if err != nil {
			erros = append(erros, "valorAporte inválido")
		}

		valorAumentoAporte, err := strconv.ParseFloat(c.Query("valorAumentoAporte"), 64)
		if err != nil {
			erros = append(erros, "valorAumentoAporte inválido")
		}

		if len(erros) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"erros": erros,
				"mensagem": "falha",
			})
			return
		}

		frequenciaAporte := financas.PegaTipoFrequencia(c.Query("frequenciaAporte"))
		frequenciaAumentoAporte := financas.PegaTipoFrequencia(c.Query("frequenciaAumentoAporte"))
		aplicacao.QuantidadeDeMeses = quantidadeDeMeses
		aplicacao.ValorInicial = valorInicial
		aplicacao.Taxa = taxa
		aporte.ValorAporte = valorAporte
		aporte.ValorAumentoAporte = valorAumentoAporte
		aporte.FrequenciaAporte = frequenciaAporte
		aporte.FrequenciaAumentoAporte = frequenciaAumentoAporte
		aplicacao.Aporte = aporte
		aplicacao.IniciarSimulacao()
		c.JSON(http.StatusOK, gin.H{
			"aplicacao": aplicacao,
			"mensagem": "sucesso",
		})
	})
}
