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
        c.log("Mensagem de chat")
        ginCtx.SSEvent(tipo, msg)
    case "entrou-chat":
        c.log("Cliente entrou em chat")
        ginCtx.SSEvent(tipo, msg)
    case "chat-nova-mensagem":
        c.log("Cliente enviou menagem ao chat")
        ginCtx.SSEvent(tipo, msg)
    case "painel":
        c.log("Mensagem de painel")
        ginCtx.SSEvent(tipo, msg)
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
        time.Sleep(5 * time.Second)
    }
}
type GerenciadorSSE struct {
    canais map[string]*CanalSSE
}
func (g *GerenciadorSSE) log(msg string) {
    dataAgora := utils.AgoraFormatado()
    fmt.Printf("[%s][SSE] >> %s\n", dataAgora, msg)
}
func (g *GerenciadorSSE) criarCanal(id string) (*CanalSSE, bool) {
    if _, ok := g.canais[id]; !ok {
        canal := CanalSSE{
            Usuario: id,
            Canal: make(chan InfoSSE),
        }
        g.canais[id] = &canal
        return g.canais[id], true
    }
    return g.canais[id], false
}
func (g *GerenciadorSSE) buscarCanal(id string) (*CanalSSE, bool) {
    if canal, existe := g.canais[id]; !existe {
        return canal, false
    }
    return g.canais[id], true
}
func (g *GerenciadorSSE) removerCanal(id string) error {
    delete(g.canais, id)
    return nil
}
type SalaChatSSE struct {
    Id string `json:"id"`
    ClientesSala []string `json:"clientes"`
}
func (s *SalaChatSSE) jaEstaEmSala(id string) bool {
    fmt.Println(s.ClientesSala)
    for _, idcliente := range s.ClientesSala {
        if idcliente == id {
            return true
        }
    }
    return false
}
func (scs *SalaChatSSE) adicionarCliente(idcliente string) {
    for _, v := range scs.ClientesSala {
        if v == idcliente {
            return
        }
    }
    scs.ClientesSala = append(scs.ClientesSala, idcliente)
}
type GerenciadorSalaChat struct {
    Salas map[string]*SalaChatSSE `json:"salas"`
}
func (gsc *GerenciadorSalaChat) criarSala(id string) *SalaChatSSE {
    if _, existe := gsc.Salas[id]; existe {
        return gsc.Salas[id]
    }
    var clientes []string
    sala := SalaChatSSE{
        Id: id,
        ClientesSala: clientes,
    }
    gsc.Salas[id] = &sala
    return &sala
}
func (gsc *GerenciadorSalaChat) buscarSala(id string) (*SalaChatSSE, bool) {
    if sala, existe := gsc.Salas[id]; existe {
        return sala, true
    }
    var sala SalaChatSSE
    return &sala, false
}
func HandlerSSE(router *gin.Engine) {
    gerenciadorSSE := GerenciadorSSE{
        canais: make(map[string]*CanalSSE),
    }
    gerenciadorSalaChat := GerenciadorSalaChat{
        Salas: make(map[string]*SalaChatSSE),
    }
    router.GET("/sse/:idusuario/entrar/:idsala", func(c *gin.Context) {
        sala := gerenciadorSalaChat.criarSala(c.Param("idsala"))
        canalSSECliente, existe := gerenciadorSSE.canais[c.Param("idusuario")]
        if !existe {
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"Cliente não tem conexao de sse aberta",
            })
            return
        }
        if sala.jaEstaEmSala(c.Param("idusuario")) {
            c.AbortWithStatus(200)
            return
        }
        sala.adicionarCliente(c.Param("idusuario"))
        var infoEntrouEmSala InfoSSE
        var conteudoInfoEntrouEmSala Conteudo
        conteudoInfoEntrouEmSala.Sala = c.Param("idsala")
        conteudoInfoEntrouEmSala.Remetente = c.Param("idusuario")
        conteudoInfoEntrouEmSala.Mensagem = "Você entrou em uma sala"
        infoEntrouEmSala.Tipo = "entrou-chat"
        infoEntrouEmSala.Conteudo = conteudoInfoEntrouEmSala
        canalSSECliente.Canal <- infoEntrouEmSala
        c.JSON(200, gin.H{
            "status":"sucesso",
            "mensagem": "sala criado e adentrada :)",
        })
        return
    })
    router.POST("/sse/:idusuario/sala/:idsala/enviar", func(c *gin.Context) {
        sala, existe := gerenciadorSalaChat.Salas[c.Param("idsala")]
        if !existe {
            c.AbortWithStatus(404)
            return
        }
        if len(c.Query("msg")) < 1 {
            c.JSON(404, gin.H{
                "status":"falha",
                "mensagem":"Conteudo de mensagem deve conter um valor!",
            })
        }
        for _, clienteid := range sala.ClientesSala {
            canalSSECliente, existe := gerenciadorSSE.canais[clienteid]
            if !existe {
                continue
            }
            var infoMensagemEnviadaACanal InfoSSE
            conteudoMnesagemEnviadaACanal := Conteudo {
                Sala: c.Param("idsala"),
                Remetente: c.Param("idusuario"),
                Mensagem: c.Query("msg"),
            }
            infoMensagemEnviadaACanal.Tipo = "chat-nova-mensagem"
            infoMensagemEnviadaACanal.Conteudo = conteudoMnesagemEnviadaACanal
            canalSSECliente.Canal <- infoMensagemEnviadaACanal
        }
        c.AbortWithStatus(201)
        return
    })
    router.GET("/sse/:idusuario", func(c *gin.Context) {
        _, existe := gerenciadorSSE.buscarCanal(c.Param("idusuario"))
        if !existe {
            canal, ok := gerenciadorSSE.criarCanal(c.Param("idusuario"))
            if !ok {
                c.AbortWithStatus(500)
            }
            fmt.Println("Canal criado >> " + canal.Usuario)
        }
        canal, _ := gerenciadorSSE.buscarCanal(c.Param("idusuario"))
        c.Writer.Header().Set("Content-Type", "text/event-stream")
        c.Writer.Header().Set("Cache-Control", "no-cache")
        c.Writer.Header().Set("Connection", "keep-alive")
        go func() {
            <-c.Request.Context().Done()
            gerenciadorSSE.removerCanal(c.Param("idusuario"))
        }()
        c.Stream(func(w io.Writer) bool {
            canal = gerenciadorSSE.canais[c.Param("idusuario")]
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
        canal, existe := gerenciadorSSE.canais[c.Param("idusuario")]
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
