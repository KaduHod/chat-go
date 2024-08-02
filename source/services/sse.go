package services

import (
	"chat/source/utils"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)
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
        //c.log("Ping")
        ginCtx.SSEvent(tipo, msg.Conteudo)
    case "chat":
        //c.log("Mensagem de chat")
        ginCtx.SSEvent(tipo, msg)
    case "entrou-chat":
        //c.log("Cliente entrou em chat")
        ginCtx.SSEvent(tipo, msg)
    case "chat-nova-mensagem":
        //c.log("Cliente enviou menagem ao chat")
        ginCtx.SSEvent(tipo, msg)
    case "painel":
        //c.log("Mensagem de painel")
        ginCtx.SSEvent(tipo, msg)
    case "deixou-sala":
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
type GerenciadorCanaisSSE struct {
    canais map[string]*CanalSSE
}
func (g *GerenciadorCanaisSSE) log(msg string) {
    dataAgora := utils.AgoraFormatado()
    fmt.Printf("[%s][SSE] >> %s\n", dataAgora, msg)
}
/*
* Criando um canal de usuario
*/
func (g *GerenciadorCanaisSSE) criarCanal(id string) (*CanalSSE, bool) {
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
// Buscando um usuario
func (g *GerenciadorCanaisSSE) buscarCanal(nomeCliente string) (*CanalSSE, bool) {
    if canal, existe := g.canais[nomeCliente]; !existe {
        return canal, false
    }
    return g.canais[nomeCliente], true
}
//removendo um usuario
func (g *GerenciadorCanaisSSE) removerCanal(id string) error {
    delete(g.canais, id)
    return nil
}
type Sala struct {
    Id string `json:"id"`
    ClientesSala []string `json:"clientes"`
}
func (s *Sala) estaEmSala(id string) bool {
    for _, idcliente := range s.ClientesSala {
        if idcliente == id {
            return true
        }
    }
    return false
}
func (s *Sala) adicionarCliente(idcliente string) {
    for _, v := range s.ClientesSala {
        if v == idcliente {
            return
        }
    }
    s.ClientesSala = append(s.ClientesSala, idcliente)
}
func (s *Sala) removerCliente(nomecliente string) {
    var novaLista []string
    for _, v := range s.ClientesSala {
        if v != nomecliente {
            novaLista = append(novaLista, v)
        }
    }
    s.ClientesSala = novaLista
}
type GerenciadorSalas struct {
    Salas map[string]*Sala `json:"salas"`
}
func (gsc *GerenciadorSalas) criarSala(id string) *Sala {
    if _, existe := gsc.Salas[id]; existe {
        return gsc.Salas[id]
    }
    var clientes []string
    sala := Sala{
        Id: id,
        ClientesSala: clientes,
    }
    gsc.Salas[id] = &sala
    return &sala
}
func (gsc *GerenciadorSalas) buscarSala(id string) (*Sala, bool) {
    if sala, existe := gsc.Salas[id]; existe {
        return sala, true
    }
    var sala Sala
    return &sala, false
}
func (gs *GerenciadorSalas) removerSala(id string) {
    if _, existe := gs.Salas[id]; existe {
        delete(gs.Salas, id)
    }
}
func HandlerSSE(router *gin.Engine) {
    gerenciadorCanais := GerenciadorCanaisSSE{
        canais: make(map[string]*CanalSSE),
    }
    gerenciadorSalas := GerenciadorSalas{
        Salas: make(map[string]*Sala),
    }
    router.GET("/chat/sse/:nomeusuario/entrar/:nomesala", func(c *gin.Context) {
        sala := gerenciadorSalas.criarSala(c.Param("nomesala"))
        _, existe := gerenciadorCanais.canais[c.Param("nomeusuario")]
        if !existe {
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"Cliente não tem conexao de sse aberta",
            })
            return
        }
        if sala.estaEmSala(c.Param("nomeusuario")) {
            c.AbortWithStatus(200)
            return
        }
        sala.adicionarCliente(c.Param("nomeusuario"))
        var infoEntrouEmSala InfoSSE
        var conteudoInfoEntrouEmSala Conteudo
        conteudoInfoEntrouEmSala.Sala = c.Param("nomesala")
        conteudoInfoEntrouEmSala.Remetente = c.Param("nomeusuario")
        conteudoInfoEntrouEmSala.Mensagem = c.Param("nomeusuario") + "entrou em uma sala"
        infoEntrouEmSala.Tipo = "entrou-chat"
        infoEntrouEmSala.Conteudo = conteudoInfoEntrouEmSala
        for _, clienteid := range sala.ClientesSala {
            canalSSE, existe := gerenciadorCanais.canais[clienteid]
            if existe {
                canalSSE.Canal <- infoEntrouEmSala
            }
        }
        c.JSON(200, gin.H{
            "status":"sucesso",
            "mensagem": "sala criado e adentrada :)",
        })
        return
    })
    router.POST("/chat/sse/:nomeusuario/sala/:nomesala/enviar", func(c *gin.Context) {
        sala, existe := gerenciadorSalas.Salas[c.Param("nomesala")]
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
        var infoMensagemEnviadaACanal InfoSSE
        conteudoMnesagemEnviadaACanal := Conteudo {
            Sala: c.Param("nomesala"),
            Remetente: c.Param("nomeusuario"),
            Mensagem: c.Query("msg"),
        }
        infoMensagemEnviadaACanal.Tipo = "chat-nova-mensagem"
        infoMensagemEnviadaACanal.Conteudo = conteudoMnesagemEnviadaACanal
        for _, clienteid := range sala.ClientesSala {
            canalSSECliente, existe := gerenciadorCanais.canais[clienteid]
            if !existe {
                continue
            }
            canalSSECliente.Canal <- infoMensagemEnviadaACanal
        }
        c.AbortWithStatus(201)
        return
    })
    router.GET("/sse/:nomeusuario", func(c *gin.Context) {
        _, existe := gerenciadorCanais.buscarCanal(c.Param("nomeusuario"))
        if !existe {
            canal, ok := gerenciadorCanais.criarCanal(c.Param("nomeusuario"))
            if !ok {
                c.AbortWithStatus(500)
            }
            fmt.Println("Canal criado >> " + canal.Usuario)
        }
        canal, _ := gerenciadorCanais.buscarCanal(c.Param("nomeusuario"))
        c.Writer.Header().Set("Content-Type", "text/event-stream")
        c.Writer.Header().Set("Cache-Control", "no-cache")
        c.Writer.Header().Set("Connection", "keep-alive")
        go func() {
            <-c.Request.Context().Done()
            gerenciadorCanais.removerCanal(c.Param("nomeusuario"))
        }()
        c.Stream(func(w io.Writer) bool {
            canal = gerenciadorCanais.canais[c.Param("nomeusuario")]
            if msg, ok := <- canal.Canal; ok {
                canal.gerenciarEventos(msg, c)
                return true
            }
            return false
        })
        return
    })
    router.GET("/chat/:nomeusuario/view", func(c *gin.Context) {
        c.HTML(200, "eventos.tmpl", gin.H{
            "title":"Eventos page",
            "nomeusuario": c.Param("nomeusuario"),
        })
        return
    })
    router.GET("/chat/sala/:nomesala/usuarios", func(c *gin.Context) {
        sala, existe := gerenciadorSalas.buscarSala(c.Param("nomesala"))
        if !existe {
            c.AbortWithStatus(404)
            return
        }
        c.JSON(200, gin.H{
            "status":"sucesso",
            "clientes": sala.ClientesSala,
        })
        return
    })
    router.GET("/chat/usuario/:nomeusuario/sala/:nomesala/sair", func(c *gin.Context) {
        canal, existe := gerenciadorCanais.buscarCanal(c.Param("nomeusuario"))
        if !existe {
            c.AbortWithStatus(404)
            return
        }
        sala, existe := gerenciadorSalas.buscarSala(c.Param("nomesala"))
        if !existe {
            c.AbortWithStatus(404)
            return
        }
        if !sala.estaEmSala(canal.Usuario) {
            c.AbortWithStatus(404)
            return
        }
        sala.removerCliente(canal.Usuario)
        if len(sala.ClientesSala) == 0 {
            gerenciadorSalas.removerSala(sala.Id)
        } else {
            var infoDeixouSala InfoSSE
            infoDeixouSala.Tipo = "deixou-sala"
            var conteudoDeixouSala Conteudo
            conteudoDeixouSala.Sala = sala.Id
            conteudoDeixouSala.Remetente = canal.Usuario
            infoDeixouSala.Conteudo = conteudoDeixouSala
            for _, nomecliente := range sala.ClientesSala {
                canalCliente, existe := gerenciadorCanais.buscarCanal(nomecliente)
                if existe {
                    canalCliente.Canal <- infoDeixouSala
                }
            }
        }
        c.AbortWithStatus(200)
        return
    })
}
