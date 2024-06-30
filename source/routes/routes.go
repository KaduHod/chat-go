package routes

import (
	"chat/source/utils"
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
	conn *websocket.Conn
	enviar chan []byte
}
func (cliente *WSCliente) ler(){
	defer func(){
		log.Println("Cliente desconectado")
		canais[canalPadrao].sair <- cliente
		cliente.conn.Close()
	}()
	for {
		_, mensagem, err := cliente.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro inexperado: %v", err)
			}
			break
		}
		canais[canalPadrao].transmissao <- mensagem
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
			cliente.conn.WriteMessage(websocket.TextMessage, mensagem)
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
	if err != nil {
		utils.Logger("/debug.log", "Erro ao iniciar websocket", "WS", false)
		log.Println(err)
		return
	}
	cliente := &WSCliente{
		conn: conn,
		enviar: make(chan []byte, 256),
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

func Router (router *gin.Engine) {
	iniciarCanalPadrao()
	router.GET("/", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message": "success",
		})
	})
	router.GET("/ws", websocketHandler)
	router.GET("/home", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message": "success",
		})		
	})
}
