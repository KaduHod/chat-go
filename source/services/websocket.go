package services

import (
	"chat/source/database"
	"chat/source/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
func (c *WSCanal) fechar(db *database.Db) error {
    fmt.Println(fmt.Sprintf("Fechando canal %s", c.nome))
    sql := fmt.Sprintf("UPDATE canal SET online = false WHERE id = %d", c.id)
    _, err := db.ExecAndLog(sql)
    if err != nil {
        return err
    }
    for cliente := range c.clientes {
        delete(c.clientes, cliente)
        fmt.Println(fmt.Sprintf("Usuario %s retirado do canal", cliente.username))
    }
    close(c.transmissao)
    close(c.entrar)
    close(c.sair)
    delete(canais, c.nome)
    return nil
}
func (c *WSCanal) registrarCanalBanco(db *database.Db) error {
    query := fmt.Sprintf("INSERT INTO canal (nome) VALUES ('%s')", c.nome)
    resultadoInsercao, err := db.ExecAndLog(query)
    if err != nil {
        return err
    }
    id, err := resultadoInsercao.LastInsertId()
    if err != nil {
        return err
    }
    c.id = id
    return nil
}
func (c *WSCanal) iniciar(db *database.Db) error {
    fmt.Println(fmt.Sprintf("Iniciando canal %s", c.nome))
    go IniciarHub(c)
    if _, jaOnline := canais[c.nome]; !jaOnline {
        canais[c.nome] = c
    }
    query := fmt.Sprintf("UPDATE canal SET online = true WHERE id = %d", c.id)
    _, err := db.ExecAndLog(query)
    if err != nil {
        if erroFecharCanal := c.fechar(db); erroFecharCanal != nil {
            fmt.Println("Erro ao inicar canal e ao tentar fechar canal")
            fmt.Println(err)
            fmt.Println(erroFecharCanal)
            return erroFecharCanal
        }
        return err
    }
    return nil
}
func (c *WSCanal) buscarCanalBanco(db *database.Db) error {
    query := fmt.Sprintf("SELECT id, nome FROM canal WHERE nome = '%s' OR id = %d LIMIT 1", c.nome, c.id)
    row := db.QueryRowAndLog(query)
    if row.Err() == sql.ErrNoRows {
        return row.Err()
    }
    if err := row.Scan(&c.id, &c.nome); err != nil {
        return err
    }
    return nil
}
func FecharCanalHandler(c *gin.Context) {
    idCanal, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "mensagem": "Falha",
            "erro": "Id de canal inválido",
        })
        return
    }
    banco := database.ConnectionConstructor()
    canal := canalConstructor(int64(idCanal), "")
    if err = canal.buscarCanalBanco(banco); err != nil {
        fmt.Println(err)
        defer banco.Conn.Close()
        c.JSON(http.StatusBadRequest, gin.H{
            "mensagem":"falha",
            "erro": "canal nao encontrado",
        })
        return
    }

    _, existe := canais[canal.nome]
    if !existe {
        c.JSON(http.StatusBadRequest, gin.H{
            "message":"falha",
            "erros": "Canal não está online para ser fechado",
        })
        return
    }
    if err = canal.fechar(banco);err != nil {
        fmt.Println("Erro ao fechar canal")
        fmt.Println(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    c.AbortWithStatus(http.StatusOK)
    return
}
func IniciarCanalHandler(c *gin.Context) {
    idCanal, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "mensagem": "Falha",
            "erro": "Id de canal inválido",
        })
        return
    }
    banco := database.ConnectionConstructor()
    canal := canalConstructor(int64(idCanal), "")
    if err := canal.buscarCanalBanco(banco); err != nil {
        defer banco.Conn.Close()
        c.JSON(http.StatusBadRequest, gin.H{
            "mensagem":"falha",
            "erro": "canal nao encontrado",
        })
        return
    }
    err = canal.iniciar(banco)
    defer banco.Conn.Close()
    if err != nil {
        fmt.Println(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    c.AbortWithStatus(http.StatusOK)
    return
}
func CriarCanalHandler(c *gin.Context) {
    nomeCanal := c.Query("nomeCanal")
    if len(nomeCanal) < 5 {
        c.JSON(http.StatusBadRequest, gin.H{
            "message":"falha",
            "erros": "nome canal inválido",
        })
        return
    }
    banco := database.ConnectionConstructor()
    // verifica se canal já está registrado
    var canal WSCanal
    canal.nome = nomeCanal
    err := canal.buscarCanalBanco(banco)
    if err != nil {
        if err = canal.registrarCanalBanco(banco); err != nil {
            fmt.Println(err)
            c.AbortWithStatus(http.StatusInternalServerError)
            defer banco.Conn.Close()
            return
        }
    }
    defer banco.Conn.Close()
    // criar canal e iniciar channels do mesmo
    canal = canalConstructor(canal.id, canal.nome)
    canais[nomeCanal]= &canal
    c.JSON(http.StatusOK, gin.H{
        "mensagem": "Canal criado!",
        "id": canal.id,
        "nome": canal.nome,
    })
    return
}
func ListarCanaisOnlineHandler(c *gin.Context) {
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
