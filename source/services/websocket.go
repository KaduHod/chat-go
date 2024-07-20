package services

import (
	"chat/source/database"
	"chat/source/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	canais = make(map[string]*WSCanal)
	canalPadrao = "canal1"
)
type WSCliente struct {
	username string
	conn *websocket.Conn
	idsocket string
	enviar chan []byte
}
type WSMensagem struct {
	remetente string
	conteudo string
}
func (cliente *WSCliente) ler(){
	defer func(){
		log.Println("Cliente desconectado")
		canais[canalPadrao].sair <- cliente
		cliente.conn.Close()
	}()
	for {
		var mensagem WSMensagem
		_, msg, err := cliente.conn.ReadMessage()

		json.Unmarshal([]byte(msg), mensagem)
		log.Println(mensagem.conteudo, mensagem.remetente)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro inexperado: %v", err)
			}
			log.Println(err)
			break
		}
		canais[canalPadrao].transmissao <- msg
		/*utils.JsonParaStruct[WSMensagem](string(msg), &mensagem)
		if mensagem.remetente != cliente.idsocket {
			canais[canalPadrao].transmissao <- []byte(mensagem.conteudo)
		}*/
	}
}
func (cliente *WSCliente) escrever(){
	defer func(){
		cliente.conn.Close()
	}()

	for {
		select {
		case mensagem, ok := <-cliente.enviar:
			if !ok {
				cliente.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := cliente.conn.WriteMessage(websocket.TextMessage, mensagem); err != nil {
				log.Println("Erro ao enviar mensagem", err)
				continue
			}
		}
	}
}
type WSCanal struct {
	clientes map[*WSCliente]bool
	transmissao chan []byte
	entrar chan *WSCliente
	sair chan *WSCliente
    nome string
    id int64
}
func (c *WSCanal) fechar() {
    fmt.Println(fmt.Sprintf("Fechando canal %s", c.nome))
    for cliente := range c.clientes {
        delete(c.clientes, cliente)
        fmt.Println(fmt.Sprintf("Usuario %s retirado do canal", cliente.username))
    }
    close(c.transmissao)
    close(c.entrar)
    close(c.sair)
}
func FecharCanalHandler(c *gin.Context) {
    var erros []string
    nomeCanal := c.Query("nomeCanal")
    if len(nomeCanal) < 5 {
        erros = append(erros, "nome canal inválido")
        c.JSON(http.StatusBadRequest, gin.H{
            "message":"falha",
            "erros": erros,
        })
        return
    }
    canal, existe := canais[nomeCanal]
    if !existe {
        erros = append(erros, "nome canal inválido")
        c.JSON(http.StatusBadRequest, gin.H{
            "message":"falha",
            "erros": erros,
        })
        return
    }
    canal.fechar()
    sql := fmt.Sprintf("UPDATE canal SET online = false WHERE id = %d", canal.id)
    banco := database.ConnectionConstructor()
    resultado, err := banco.ExecAndLog(sql)
    defer banco.Conn.Close()
    fmt.Println(resultado)
    if err != nil {
        fmt.Println("Erro ao deletar canal de banco")
        fmt.Println(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    c.AbortWithStatus(http.StatusOK)
    return
}
func IniciarCanalHandler(c *gin.Context) {

}
func CriarCanalHandler(c *gin.Context) {
    var erros []string
//    var _canal WSCanal
    //var id int64
    nomeCanal := c.Query("nomeCanal")
    if len(nomeCanal) < 5 {
        erros = append(erros, "nome canal inválido")
        c.JSON(http.StatusBadRequest, gin.H{
            "message":"falha",
            "erros": erros,
        })
        return
    }
    banco := database.ConnectionConstructor()
    // verifica se canal já está registrado
    query := fmt.Sprintf("SELECT id FROM canal WHERE nome = '%s' LIMIT 1", nomeCanal)
    linha := banco.QueryRowAndLog(query)
    if linha.Err() == sql.ErrNoRows {
        //criar registro no banco
        query = fmt.Sprintf("INSERT INTO canal (nome) VALUES ('%s')", nomeCanal)
        result, err := banco.ExecAndLog(query)
        defer banco.Conn.Close()
        if err != nil {
            fmt.Println(err)
            if utils.VerificaPadrao("Duplicate entry '([^']*)' for key 'canal.nome'", err.Error()) {
                c.JSON(http.StatusBadRequest, gin.H {
                    "mensagem": "falha",
                    "erro" : fmt.Sprintf("Nome de canal %s já está sendo utilizado", nomeCanal),
                })
                return
            }
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        _, err = result.LastInsertId()
        if err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
    } else {
    //    var canal WSCanal
    }
    /*

    canal := canalConstructor(id, nomeCanal)
    canais[nomeCanal]= &canal
    c.JSON(http.StatusOK, gin.H{
        "mensagem": "Canal criado!",
        "id": id,
        "nome": nomeCanal,
    })
    return*/
}
func ListarCanaisHandler(c *gin.Context) {
    listaCanais := make([]string, 0, len(canais))
    for key:= range canais {
        listaCanais = append(listaCanais, key)
    }

    c.JSON(http.StatusOK, gin.H{
        "canais": listaCanais,
        "mensagem": "sucesso",
    })
}
func WebsocketHandler(c * gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	idConexaoSocket := c.Request.Header.Get("Sec-WebSocket-Key")
	fmt.Println(idConexaoSocket)
	if err != nil {
		utils.Logger("/debug.log", "Erro ao iniciar websocket", "WS", false)
		log.Println(err)
		return
	}
	username := c.Param("Username")
	//secretKey := c.Param("tk")
	cliente := &WSCliente{
		conn: conn,
		idsocket: idConexaoSocket,
		enviar: make(chan []byte, 256),
		username: username,
	}
	canalPadrao := "canal1"
	canais[canalPadrao].entrar <- cliente
	go cliente.escrever()
	go cliente.ler()
}
func canalConstructor(id int64, nome string) WSCanal {
    canal := &WSCanal{
        clientes:    make(map[*WSCliente]bool),
        transmissao: make(chan []byte),
        entrar:      make(chan *WSCliente),
        sair:        make(chan *WSCliente),
        id: id,
        nome: nome,
    }
    return *canal
}
func IniciarCanalPadrao() {
	defaultChannel := &WSCanal{
		clientes:    make(map[*WSCliente]bool),
		transmissao: make(chan []byte),
		entrar:      make(chan *WSCliente),
		sair:        make(chan *WSCliente),
	}
	canais[canalPadrao] = defaultChannel
	log.Println("Canal padrão inicializado")
	go IniciarHub(defaultChannel)
}

func IniciarHub(canal *WSCanal) {
	for {
		select {
		case cliente := <-canal.entrar:
			log.Println("Cliente entrou no canal")
			canal.clientes[cliente] = true
		case cliente := <-canal.sair:
			if _, ok := canal.clientes[cliente]; ok {
				log.Println("Cliente saiu do canal")
				delete(canal.clientes, cliente)
				close(cliente.enviar)
			}
		case mensagem := <-canal.transmissao:
			log.Printf("Transmitindo mensagem: %s", mensagem)
			for cliente := range canal.clientes {
				fmt.Println("Mensagem >>", cliente.idsocket)
				select {
				case cliente.enviar <- mensagem:
				default:
					close(cliente.enviar)
					delete(canal.clientes, cliente)
				}
			}
		}
	}
}
