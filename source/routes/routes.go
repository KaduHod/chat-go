package routes

import (
	"chat/source/database"
	"chat/source/utils"
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
		//var mensagem WSMensagem
		_, msg, err := cliente.conn.ReadMessage()
		log.Println(string(msg))
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
			/*msg := WSMensagem{
				remetente: cliente.idsocket,
				conteudo: string(mensagem),
			}
			if err := cliente.conn.WriteMessage(websocket.TextMessage, []byte(utils.StructParaJson[WSMensagem](msg))); err != nil {
				log.Println("Erro ao escrever json")
				log.Println(err)
				continue
			}*/
		}
	}
}
type WSCanal struct {
	clientes map[*WSCliente]bool
	transmissao chan []byte
	entrar chan *WSCliente
	sair chan *WSCliente
}
func websocketHandler(c * gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	idConexaoSocket := c.Request.Header.Get("Sec-WebSocket-Key")
	fmt.Println(idConexaoSocket)
	if err != nil {
		utils.Logger("/debug.log", "Erro ao iniciar websocket", "WS", false)
		log.Println(err)
		return
	}
	username := c.Param("Username")
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
func iniciarCanalPadrao() {
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
type Usuario struct {
	Id int `json:"id"`
	Nome string `json:"nome"`
	Apelido string `json:"apelido"`
}
func Router (router *gin.Engine) {
	iniciarCanalPadrao()
	router.GET("/", func (c *gin.Context) {
		c.Header("Content-type", "text/html")
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.POST("/logar", func(c *gin.Context) {
		var usr Usuario
		if !utils.RequestBody[Usuario](&usr, c) {
			return
		} 
		sql := fmt.Sprintf("INSERT INTO usuario (nome, apelido) VALUES ('%s','%s')", usr.Nome, usr.Apelido)
		db := database.ConnectionConstructor()
		db.ExecAndLog(sql)
		sql = fmt.Sprintf("SELECT id, nome, apelido FROM usuario WHERE apelido = '%s' LIMIT 1", usr.Apelido)
		row := db.QueryRowAndLog(sql)
		if row.Err() != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			defer db.Conn.Close()
			return
		}
		if err := row.Scan(&usr.Id, &usr.Nome, &usr.Apelido); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			defer db.Conn.Close()
			return
		}
		defer db.Conn.Close()
		c.JSON(http.StatusOK, gin.H{
			"mensagem": "ok",
			"usuario": usr,
			"chave_secreta": utils.HashSha512(fmt.Sprintf("%d", usr.Id)),
		})
		return 
	})
	router.GET("/ws", websocketHandler)
}
