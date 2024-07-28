package services

import (
	"chat/source/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)
const (
    CANAIS_DN="/app/canais.dn"
    LIMITE_CANAIS_EM_MEMORIA=10
    LIMITE_CLIENTES_EM_MEMORIA=100
)
type EventStreamRequest struct {
    Mensagem string `form:"mensagem" json:"mensagem" binding:"required,max=100"`
}
type DaemonSSE struct {
    listaCanaisConf ListaCanalConf
    listaCanaisEmMemoria map[string]*Canal
    listaClientesEmMemoria map[string]*Cliente
}
func (d *DaemonSSE) lerArquivoConf() {
    if err := utils.LerArquivoJson[ListaCanalConf](CANAIS_DN, &d.listaCanaisConf); err != nil {
        fmt.Println("Erro ao ler arquivo de configuração do daemon")
        panic(err)
    }
}
func (d *DaemonSSE) pegaCanalEmMemoria(id string) (*Canal, bool) {
    canal, ok := d.listaCanaisEmMemoria[id]
    if ok {
        return canal, ok
    }
    return nil, false
}
func (d *DaemonSSE) pegaClienteEmMemoria(id string) (*Cliente, bool) {
    cliente, ok := d.listaClientesEmMemoria[id]
    if ok {
        return cliente, ok
    }
    return nil, false
}
func (d *DaemonSSE) loop() {
    for {
        d.lerArquivoConf()
        for idcanal, canalConf := range d.listaCanaisConf.Lista {
            d.log(fmt.Sprintf("Canal -> %s", idcanal))
            for _, cliente := range canalConf.Clientes {
                d.log(fmt.Sprintf("\tCliente -> %s", cliente))
            }
        }
        time.Sleep(4 * time.Second)
    }
}
func (s *DaemonSSE) log(conteudo string) {
    agora := utils.AgoraFormatado()
    log := fmt.Sprintf("[%s][DAEMON SSE] >> %s",agora ,conteudo)
    fmt.Println(log)
}
type ListaCanalConf struct {
    Lista map[string]*CanalConf `json:"lista"`
}
func (conf *ListaCanalConf) adicionaCanal(idcanal string) {
    var listaClientes []string
    conf.Lista[idcanal] = &CanalConf{
        Clientes: listaClientes,
    }
}
func (conf *ListaCanalConf) adicionaClienteCanal(idcanal string ,idcliente string) error {
    if _, existe := conf.Lista[idcanal]; !existe {
        return errors.New("Canal nao existe para adicionar usuario")
    }
    conf.Lista[idcanal].Clientes = append(conf.Lista[idcanal].Clientes, idcliente)
    return nil
}
func (conf *ListaCanalConf) atualizarArquivoConf() {
    json, err := json.Marshal(conf)
    if err != nil {
        panic(err)
    }
    if err := utils.SobrescreverArquivo(CANAIS_DN, string(json)); err != nil {
        panic(err)
    }
}
type CanalConf struct {
    Clientes []string `json:"clientes"`
}
type Canal struct {
    id string
    canal chan string
}
type Cliente struct {
    id string
    contextGin *gin.Context
}
func (d *DaemonSSE) criarCanal(id string) (*Canal) {
    var canal Canal
    if canalExistente, existe := d.pegaCanalEmMemoria(id); existe {
        return canalExistente
    }
    d.listaCanaisEmMemoria[id] = &canal
    return &canal
}
func ControladorServerSendEvents(router *gin.Engine) {
    var listaCanalConf ListaCanalConf
    listaCanalConf.Lista = map[string]*CanalConf{}
    daemonSSE := DaemonSSE{
        listaCanaisConf: listaCanalConf,
        listaCanaisEmMemoria: make(map[string]*Canal),
        listaClientesEmMemoria: make(map[string]*Cliente),
    }
    utils.EscreverEmArquivo(CANAIS_DN, "{}")
    go func() {
        daemonSSE.loop()
    }()
    router.GET("/sse/usuario/:idusuario/canal/:idcanal", func (c *gin.Context) {
        iduser := c.Param("idusuario")
        idcanal := c.Param("idcanal")
        if _, existe := daemonSSE.pegaCanalEmMemoria(idcanal); !existe {
            fmt.Println("Criando canal")
            canal := &Canal{
                id: idcanal,
                canal: make(chan string),
            }
            daemonSSE.listaCanaisEmMemoria[idcanal] = canal
        }
        if _, existe := daemonSSE.pegaClienteEmMemoria(iduser); !existe {
            fmt.Println("Criando cliente")
            cliente := &Cliente{
                id: iduser,
                contextGin: c,
            }
            daemonSSE.listaClientesEmMemoria[iduser] = cliente
        }
        daemonSSE.lerArquivoConf()
        _, existe := daemonSSE.listaCanaisConf.Lista[idcanal]
        if !existe {
            daemonSSE.listaCanaisConf.adicionaCanal(idcanal)
        }
        if err := daemonSSE.listaCanaisConf.adicionaClienteCanal(idcanal, iduser); err != nil {
            fmt.Println(err)
            fmt.Printf("Erro canal: %s cliente: %s", idcanal, iduser)
            delete(daemonSSE.listaClientesEmMemoria, iduser)
            daemonSSE.listaCanaisConf.atualizarArquivoConf()
            c.AbortWithStatus(400)
            return
        }
        daemonSSE.listaCanaisConf.atualizarArquivoConf()
        return
    })
    //var gerenciadorCanais GerenciadorCanais
   // canal := make(chan string)
   /*
    router.POST("/event-stream", func (c *gin.Context) {
        canal, ok := gerenciadorCanais.canais[c.Param("idcanal")]
        if !ok {
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"Canal nao encontrado",
            })
            return
        }
        ControladorEventosPost(c, canal)
    })
    router.GET("/event-stream/:idcanal", func (c *gin.Context) {
        canal, ok := gerenciadorCanais.canais[c.Param("idcanal")]
        if !ok {
            canal = gerenciadorCanais.criarCanal(c.Param("idncanal"))
        }
        canal.adicionaConexao(c)
        //ControladorEventosGet(c, canal)
    })
    */
}
func ControladorEventosPost(c *gin.Context, canal Canal) {
    var requisicao EventStreamRequest
    if err := c.ShouldBind(&requisicao); err != nil {
        c.JSON(400, gin.H{
            "status": "falha",
            "mensagem": err.Error(),
        })
        return
    }
    //fmt.Println(requisicao, canal)
    canal.canal <- requisicao.Mensagem
    c.JSON(201, gin.H{
        "status":"criado",
    })
    return
}
func ControladorEventosGet(c *gin.Context, ch chan string) {
    c.Stream(func(w io.Writer) bool {
        if msg, ok := <-ch; ok {
            c.SSEvent("message", msg)
            return true
        }
        return false
    })

    return
}

