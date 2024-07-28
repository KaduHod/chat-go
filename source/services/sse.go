package services

import (
	"chat/source/utils"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)
type EventStreamRequest struct {
    Mensagem string `form:"mensagem" json:"mensagem" binding:"required,max=100"`
}
type Conteudo struct {
    Mensagem string `json:"mensagem"`
    Remetente string `json:"remetente"`
    Sala string `json:"sala"`
}
type InfoSSE struct {
    Tipo string `json:"tipo" binding:"required"`
    Conteudo Conteudo `json:"conteudo" binding:"required"`
}
type CanalSSE struct {
    Usuario string `json:"usuario"`
    Canal chan InfoSSE
}
func (c *CanalSSE) gerenciarEventos(msg InfoSSE, ginCtx *gin.Context) {
    switch tipo := msg.Tipo; tipo {
    case "ping":
        c.log("Ping")
        ginCtx.SSEvent(tipo, msg.Conteudo)
    case "chat":
        ginCtx.SSEvent(tipo, msg.Conteudo)
        c.log("Mensagem de chat")
    case "painel":
        c.log("Mensagem de painel")
    }
}
func (c *CanalSSE) log(msg string) {
    dataAgora := utils.AgoraFormatado()
    fmt.Printf("[%s][SSE] >> %s\n", dataAgora, msg)
}
func (c *CanalSSE) ping() {
    for {
        conteudo := Conteudo {
            Mensagem: "ping",
            Remetente: "Sistema",
        }
        c.Canal <- InfoSSE{
            Tipo: "ping",
            Conteudo: conteudo,
        }
        time.Sleep(10 * time.Second)
    }
}
type GerenciadorSSE struct {
    canais map[string]CanalSSE
}
func (g *GerenciadorSSE) log(msg string) {
    dataAgora := utils.AgoraFormatado()
    fmt.Printf("[%s][SSE] >> %s\n", dataAgora, msg)
}
func (g *GerenciadorSSE) criarCanal(id string) (CanalSSE, bool) {
    if _, ok := g.canais[id]; !ok {
        canal := CanalSSE{
            Usuario: id,
            Canal: make(chan InfoSSE),
        }
        g.canais[id] = canal
        return canal, true
    }
    return g.canais[id], false
}
func (g *GerenciadorSSE) buscarCanal(id string) (CanalSSE, bool) {
    var canal CanalSSE
    if canal, existe := g.canais[id]; !existe {
        return canal, false
    }
    return canal, true
}
func HandlerSSE(router *gin.Engine) {
    gerenciador := GerenciadorSSE{
        canais: make(map[string]CanalSSE),
    }
    router.GET("/sse/:idusuario", func(c *gin.Context) {
        canal, existe := gerenciador.buscarCanal(c.Param("idusuario"))
        if !existe {
            canal, ok := gerenciador.criarCanal(c.Param("idusuario"))
            fmt.Println(canal)
            if !ok {
                c.AbortWithStatus(500)
            }
            fmt.Println("Canal criado >> " + canal.Usuario)
        }
        c.Writer.Header().Set("Content-Type", "text/event-stream")
        c.Writer.Header().Set("Cache-Control", "no-cache")
        c.Writer.Header().Set("Connection", "keep-alive")
        go canal.ping()
        c.Stream(func(w io.Writer) bool {
            canal = gerenciador.canais[c.Param("idusuario")]
            if msg, ok := <- canal.Canal; ok {
                canal.gerenciarEventos(msg, c)
                return true
            }
            return false
        })
        return
    })
    router.POST("/sse/:idusuario", func(c *gin.Context) {
        var infoParaPostar InfoSSE
        if err := c.ShouldBind(&infoParaPostar); err != nil {
            fmt.Println("Erro >> ", err)
            c.AbortWithStatus(400)
            return
        }
        canal, existe := gerenciador.canais[c.Param("idusuario")]
        if !existe {
            c.AbortWithStatus(404)
            return
        }
        canal.Canal <- infoParaPostar
        c.JSON(201, gin.H{
            "status":"sucesso",
        })
        return
    })
    router.GET("/sse/:idusuario/view", func(c *gin.Context) {
        c.HTML(200, "eventos.tmpl", gin.H{
            "title":"Eventos page",
            "idusuario": c.Param("idusuario"),
        })
        return
    })
}
