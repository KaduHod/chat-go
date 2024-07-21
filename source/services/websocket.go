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
type WSCanalCliente struct {
    id int64
    idUsuario int64
    idCanal int64
    online bool
}
func (cu *WSCanalCliente) insereRegristro(db *database.Db) error {
    query := fmt.Sprintf("INSERT INTO usuariocanal (idusuario, idcanal, online) VALUES (%d, %d, true)", cu.idUsuario, cu.idCanal)
    _, err := db.ExecAndLog(query)
    if err != nil {
        fmt.Println(err)
        return err
    }
    return nil
}
func (cu *WSCanalCliente) buscaRegistro(db *database.Db) error {
    query := fmt.Sprintf("SELECT id, online FROM usuariocanal WHERE idusuario = %d AND idcanal = %d", cu.idUsuario, cu.idCanal)
    linha := db.QueryRowAndLog(query)
    if linha.Err() != nil {
        return linha.Err()
    }
    if err := linha.Scan(&cu.id, &cu.online); err != nil {
        return err
    }
    return nil
}
func (cu *WSCanalCliente) alteraParaOnline(db *database.Db) error {
    query := fmt.Sprintf("UPDATE usuariocanal SET online = true WHERE id = %d", cu.id)
    _, err := db.ExecAndLog(query)
    if err != nil {
        fmt.Println(err)
        return err
    }
    cu.online = true
    return nil
}
func (cu *WSCanalCliente) alteraParaOffline(db *database.Db) error {
    query := fmt.Sprintf("UPDATE usuariocanal SET online = false WHERE id = %d", cu.id)
    _, err := db.ExecAndLog(query)
    if err != nil {
        fmt.Println(err)
        return err
    }
    cu.online = false
    return nil
}
type WSCliente struct {
    username string
    id int64
    conn *websocket.Conn
    idsocket string
    enviar chan *WSMensagem
}
func (c *WSCliente) buscarRegistro(db *database.Db) error {
    query := fmt.Sprintf("SELECT apelido FROM usuario WHERE id = %d", c.id)
    linha := db.QueryRowAndLog(query)
    if linha.Err() != nil {
        return linha.Err()
    }
    if err := linha.Scan(&c.username); err != nil {
        return err
    }
    return nil
}
type WSMensagem struct {
    Remetente string `json:"remetente"`
    Conteudo string `json:"conteudo"`
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
        json.Unmarshal([]byte(msg), &mensagem)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Erro inexperado: %v", err)
            }
            log.Println(err)
            break
        }
        canais[canalPadrao].transmissao <- &mensagem
    }
}
func (cliente *WSCliente) lerV2(nomeCanal string){
    defer func(){
        log.Println("Cliente desconectado")
        canais[nomeCanal].sair <- cliente
        cliente.conn.Close()
    }()
    for {
        var mensagem WSMensagem
        _, msg, err := cliente.conn.ReadMessage()

        json.Unmarshal([]byte(msg), &mensagem)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Erro inexperado: %v", err)
            }
            log.Println(err)
            break
        }
        mensagem.Conteudo = string(msg)
        mensagem.Remetente = cliente.username
        canais[nomeCanal].transmissao <- &mensagem
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
            if err := cliente.conn.WriteJSON(mensagem); err != nil {
                log.Println("Erro ao enviar mensagem", err)
                continue
            }
        }
    }
}
func (cliente *WSCliente) escreverV2(){
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
            if err:= cliente.conn.WriteJSON(mensagem); err != nil {
                log.Println("Erro ao enviar mensagem", err)
                continue
            }
        }
    }
}
type WSCanal struct {
    clientes map[*WSCliente]bool
    transmissao chan *WSMensagem
    entrar chan *WSCliente
    sair chan *WSCliente
    nome string
    id int64
}
func (c *WSCanal) removerUsuario(cliente *WSCliente) error {
    canalUsuario := WSCanalCliente{
        idUsuario: cliente.id,
        idCanal: c.id,
    }
    banco := database.ConnectionConstructor()
    if err := canalUsuario.buscaRegistro(banco); err != nil {
        defer banco.Conn.Close()
        fmt.Println(err)
        return err
    }
    if err := canalUsuario.alteraParaOffline(banco); err != nil {
        defer banco.Conn.Close()
        fmt.Println(err)
        return err
    }
    defer banco.Conn.Close()
    delete(c.clientes, cliente)
    close(cliente.enviar)
    return nil
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
func (c *WSCanal) buscarRegistro(db *database.Db) error {
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
type WSCanalUsuarioOnline struct {
    Apelido string `json:"username"`
    Online bool `json:"online"`
}
func (c *WSCanal) buscarClientes(db *database.Db) ([]WSCanalUsuarioOnline, error) {
    var clientes []WSCanalUsuarioOnline
    query := fmt.Sprintf("SELECT u.apelido, uc.online FROM usuario u INNER JOIN usuariocanal uc on uc.idusuario = u.id AND uc.idcanal = %d", c.id)
    linhas, err := db.QueryAndLog(query)
    if err != nil {
        return clientes, err
    }
    for linhas.Next() {
        var usuario WSCanalUsuarioOnline
        if err = linhas.Scan(&usuario.Apelido, &usuario.Online); err != nil {
            return clientes, err
        }
        clientes = append(clientes, usuario)
    }
    return clientes, nil
}
func FecharCanalHandler(c *gin.Context) {
    idCanal, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "Falha",
            "erro": "Id de canal inválido",
        })
        return
    }
    banco := database.ConnectionConstructor()
    canal := canalConstructor(int64(idCanal), "")
    if err = canal.buscarRegistro(banco); err != nil {
        fmt.Println(err)
        defer banco.Conn.Close()
        c.JSON(http.StatusBadRequest, gin.H{
            "status":"falha",
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
            "status": "Falha",
            "erro": "Id de canal inválido",
        })
        return
    }
    banco := database.ConnectionConstructor()
    canal := canalConstructor(int64(idCanal), "")
    if err := canal.buscarRegistro(banco); err != nil {
        defer banco.Conn.Close()
        c.JSON(http.StatusBadRequest, gin.H{
            "status":"falha",
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
    err := canal.buscarRegistro(banco)
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
        "status": "Canal criado!",
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
        "status": "sucesso",
    })
}
func AdicionarUsuarioCanalHandler(c *gin.Context) {
    userId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H {
            "status": "falha",
            "erro": "id de usuario invalido",
        })
        return
    }
    canalId, err := strconv.Atoi(c.Param("idcanal"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H {
            "status": "falha",
            "erro": "id de canal invalido",
        })
        return
    }
    cliente := clienteConstructor(int64(userId))
    banco := database.ConnectionConstructor()
    defer banco.Conn.Close()
    if err = cliente.buscarRegistro(banco); err != nil {
        defer banco.Conn.Close()
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{
                "status": "falha",
                "erro": "Usuario nao encontrado",
            })
            return
        }
        fmt.Println(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    canal := canalConstructor(int64(canalId), "")
    if err = canal.buscarRegistro(banco); err != nil {
        defer banco.Conn.Close()
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{
                "status": "falha",
                "erro": "Canal nao encontrado",
            })
            return
        }
        fmt.Println(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    canalUsuario := WSCanalCliente{
        idUsuario: cliente.id,
        idCanal: canal.id,
    }
    err = canalUsuario.buscaRegistro(banco)
    if err != sql.ErrNoRows && err != nil {
        fmt.Println(err)
        defer banco.Conn.Close()
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    //cria registro
    if err == sql.ErrNoRows {
        canalUsuario.online = true
        if err = canalUsuario.insereRegristro(banco); err != nil {
            defer banco.Conn.Close()
            fmt.Println(err)
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        } else {
            defer banco.Conn.Close()
            c.AbortWithStatus(http.StatusOK)
            return
        }
    }
    if canalUsuario.online {
        defer banco.Conn.Close()
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "falha",
            "erro": "Usuario já está no canal!",
        })
        return
    }
    //atualiza registro
    if !canalUsuario.online {
        if err = canalUsuario.alteraParaOnline(banco); err != nil {
            fmt.Println(err)
            defer banco.Conn.Close()
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
    }
    defer banco.Conn.Close()
    c.AbortWithStatus(http.StatusOK)
    return
}
func ListaUsuariosDeCanalHandler(c *gin.Context) {
    idCanal, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(404, gin.H{
            "erro": "Id de canal inválido",
            "status": "falha",
        })
        return
    }
    canal := canalConstructor(int64(idCanal), "")
    banco := database.ConnectionConstructor()
    if err = canal.buscarRegistro(banco); err != nil {
        defer banco.Conn.Close()
        if err == sql.ErrNoRows {
            c.JSON(200, gin.H{
                "status": "falha",
                "erro": "canal nao encontrado",
            })
            return
        }
        c.AbortWithStatus(500)
        return
    }
    clientes, err := canal.buscarClientes(banco)
    defer banco.Conn.Close()
    if err != nil {
        c.AbortWithStatus(500)
        return
    }
    c.JSON(200, gin.H{
        "status": "sucesso",
        "clientes": clientes,
    })
    return
}
func RemoverUsuarioCanalHandler(c *gin.Context) {
    userId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H {
            "status": "falha",
            "erro": "id de usuario invalido",
        })
        return
    }
    canalId, err := strconv.Atoi(c.Param("idcanal"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H {
            "status": "falha",
            "erro": "id de canal invalido",
        })
        return
    }
    banco := database.ConnectionConstructor()
    canalUsuario := WSCanalCliente{
        idUsuario: int64(userId),
        idCanal: int64(canalId),
    }
    if err = canalUsuario.buscaRegistro(banco); err != nil {
        defer banco.Conn.Close()
        if err == sql.ErrNoRows {
            c.JSON(http.StatusBadRequest, gin.H{
                "status": "falha",
                "erro": "Usuario não está no canal!",
            })
        }
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    if !canalUsuario.online {
        defer banco.Conn.Close()
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "falha",
            "erro": "Usuário nao está online no canal",
        })
        return
    }
    err = canalUsuario.alteraParaOffline(banco)
    defer banco.Conn.Close()
    if  err != nil {
        fmt.Println(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    c.AbortWithStatus(http.StatusOK)
    return
}
func WebsocketHandlerV2(c *gin.Context) {
    idCanal, err := strconv.Atoi(c.Param("idcanal"))
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "falha",
            "erro": "Id de canal nao inválido!",
        })
        return
    }
    idUsuario, err := strconv.Atoi(c.Param("idusuario"))
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "falha",
            "erro": "Id de usuario nao inválido!",
        })
        return
    }
    banco := database.ConnectionConstructor()
    canal := canalConstructor(int64(idCanal), "")
    if err = canal.buscarRegistro(banco); err != nil {
        defer banco.Conn.Close()
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "falha",
            "erro": "Canal não encontrado",
        })
        return
    }
    cliente := clienteConstructor(int64(idUsuario))
    if err = cliente.buscarRegistro(banco); err != nil {
        defer banco.Conn.Close()
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "falha",
            "erro": "Usuário não encontrado",
        })
        return
    }
    canalUsuario := WSCanalCliente{
        idUsuario: cliente.id,
        idCanal: canal.id,
    }
    err = canalUsuario.buscaRegistro(banco)
    if err != sql.ErrNoRows && err != nil {
        fmt.Println(err)
        defer banco.Conn.Close()
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    //cria registro
    if err == sql.ErrNoRows {
        canalUsuario.online = true
        if err = canalUsuario.insereRegristro(banco); err != nil {
            defer banco.Conn.Close()
            fmt.Println(err)
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
    }
    if canalUsuario.online {
        defer banco.Conn.Close()
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "falha",
            "erro": "Usuario já está no canal!",
        })
        return
    }
    //atualiza registro
    if !canalUsuario.online {
        if err = canalUsuario.alteraParaOnline(banco); err != nil {
            fmt.Println(err)
            defer banco.Conn.Close()
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
    }
    _, existe := canais[canal.nome]
    if !existe {
        canais[canal.nome] = &canal
        fmt.Println("Iniciando canal "+canal.nome)
        go IniciarHub(canais[canal.nome])
    }
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        utils.Logger("/debug.log", "Erro ao iniciar websocket", "WS", false)
        log.Println(err)
        return
    }
    idConexaoSocket := c.Request.Header.Get("Sec-WebSocket-Key")
    cliente.conn = conn
    cliente.idsocket = idConexaoSocket
    cliente.enviar = make(chan *WSMensagem, 256)
    canais[canal.nome].entrar <- &cliente
    go cliente.escreverV2()
    go cliente.lerV2(canal.nome)
    c.AbortWithStatus(200)
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
    cliente := &WSCliente{
        conn: conn,
        idsocket: idConexaoSocket,
        enviar: make(chan *WSMensagem, 256),
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
        transmissao: make(chan *WSMensagem),
        entrar:      make(chan *WSCliente),
        sair:        make(chan *WSCliente),
        id: id,
        nome: nome,
    }
    return *canal
}
func clienteConstructor(id int64) WSCliente {
    cliente := WSCliente{
        id: id,
    }
    return cliente
}
func IniciarCanalPadrao() {
    defaultChannel := &WSCanal{
        clientes:    make(map[*WSCliente]bool),
        transmissao: make(chan *WSMensagem),
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
                fmt.Println(cliente.username)
                if err := canal.removerUsuario(cliente); err != nil {
                    delete(canal.clientes, cliente)
                    close(cliente.enviar)
                }
            }
        case mensagem := <-canal.transmissao:
            for cliente := range canal.clientes {
                select {
                case cliente.enviar <- mensagem:
                default:
                    if err := canal.removerUsuario(cliente); err != nil {
                        delete(canal.clientes, cliente)
                        close(cliente.enviar)
                    }
                }
            }
        }
    }
}
