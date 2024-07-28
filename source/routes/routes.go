package routes

import (
	"chat/source/database"
	"chat/source/modules/financas"
	"chat/source/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)
func Router (router *gin.Engine) {
	router.Static("/public", "./public")
    router.LoadHTMLGlob("templates/*")
    router.GET("/", func (c *gin.Context) {
        c.JSON(200, gin.H{
            "mensagem": "Olá, algum erro aconteceu para voce estar aqui!",
        })
    })
    router.GET("/home/:nome", func (c *gin.Context) {
		c.Header("Content-type", "text/html")
        nomeUsuario := c.Param("nome")
        usuario := services.WSCliente{
            Username: nomeUsuario,
        }
        banco := database.ConnectionConstructor()
        err := usuario.BuscarRegistro(banco)
        fmt.Println(usuario)
        defer banco.Conn.Close()
        if err != nil {
            println(err)
            c.Redirect(302, "/")
            return
        }
        listaCanais, err := services.ListarCanaisService()
        if err != nil {
            c.HTML(http.StatusOK, "index.tmpl", gin.H{
                "titulo": "Tivemos um problema :(. Volte mais tarde.",
            })
            return
        }
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "titulo": "Bem vindo!",
            "canais": listaCanais,
            "usuario": usuario,
        })
	})
    router.GET("/home/canal/:id/usuario/:nome", func(c *gin.Context){
        idCanal, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.Redirect(302, "/")
            return
        }
        nomeUsuario := c.Param("nome")
        usuario := services.WSCliente{
            Username: nomeUsuario,
        }
        banco := database.ConnectionConstructor()
        if err = usuario.BuscarRegistro(banco); err != nil {
            println(err)
            defer banco.Conn.Close()
            c.Redirect(302, "/")
            return
        }
        canal := services.WSCanal{
            Id: int64(idCanal),
        }
        err = canal.BuscarRegistro(banco)
        defer banco.Conn.Close()
        if err != nil {
            println(err)
            c.Redirect(302, "/")
            return
        }
        fmt.Println(usuario)
        c.HTML(200, "chatV2.tmpl", gin.H{
            "canal": canal,
            "usuario": usuario,
        })
    })
	router.GET("/home/chat", func (c *gin.Context) {
		c.Header("Content-type", "text/html")
		c.HTML(http.StatusOK, "chat.tmpl", gin.H{
            "titulo": "Ola mundo",
        })
	})
	//services.IniciarCanalPadrao()
	router.GET("/financas/grafico", func (c* gin.Context) {
		c.HTML(http.StatusOK, "grafico.tmpl", gin.H{})
	})
	router.POST("/logar", services.LogarUsuario)
    router.GET("/chat/canais/online", services.ListarCanaisOnlineHandler)
    router.GET("/chat/canais", services.ListarCanaisHandler)
    router.GET("/chat/criar", services.CriarCanalHandler)
    router.GET("/chat/canal/:id/iniciar", services.IniciarCanalHandler)
    router.GET("/chat/canal/:id/fechar", services.FecharCanalHandler)
    router.GET("/chat/canal/:id/usuarios", services.ListaUsuariosDeCanalHandler)
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
    services.HandlerSSE(router)
    router.GET("/eventos/teste", func(c *gin.Context) {
        c.HTML(200, "eventos.tmpl", gin.H{})
    })
}
