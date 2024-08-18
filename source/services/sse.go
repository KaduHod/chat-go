package services

import (
	"chat/source/database"
	"chat/source/utils"
	"database/sql"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)
type Conteudo struct {
    Mensagem string `json:"mensagem"`
    Remetente string `json:"remetente"`
    Sala string `json:"sala"`
}
func newConteudo(mensagem string, apelidoUsuario string, nomesala string) *Conteudo {
    return &Conteudo{
        Mensagem: mensagem,
        Remetente: apelidoUsuario,
        Sala: nomesala,
    }
}
type InfoSSE struct {
    Tipo string `json:"tipo" binding:"required"`
    Conteudo *Conteudo `json:"conteudo" binding:"required"`
}
func newInfoSSE(tipo string,conteudo *Conteudo) *InfoSSE {
    return &InfoSSE{
        Conteudo: conteudo,
        Tipo: tipo,
    }
}
type UsuarioBD struct {
    Id int64 `json:"id"`
    Nome string `form:"nome" json:"nome" bind:"required"`
    Apelido string `form:"apelido" json:"apelido" bind:"required"`
}
type CanalSSE struct {
    UsuarioBd UsuarioBD
    Canal chan *InfoSSE
}
func (c *CanalSSE) gerenciarEventos(msg *InfoSSE, ginCtx *gin.Context) {
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
    info := newInfoSSE("ping", newConteudo("ping", "sistema", ""))
    for {
        c.Canal <- info
        time.Sleep(5 * time.Second)
    }
}
type GerenciadorCanaisSSE struct {
    canais sync.Map
    //canais map[string]*CanalSSE
}
func (g *GerenciadorCanaisSSE) log(msg string) {
    dataAgora := utils.AgoraFormatado()
    fmt.Printf("[%s][SSE] >> %s\n", dataAgora, msg)
}
/*
* Criando um canal de usuario
*/
func (g *GerenciadorCanaisSSE) criarCanal(usuarioBd UsuarioBD) (*CanalSSE, bool) {
    if canal, existe := g.buscarCanal(usuarioBd.Apelido); existe {
        return canal, true
    } else {
        canal := CanalSSE{
            UsuarioBd: usuarioBd,
            Canal: make(chan *InfoSSE),
        }
        g.canais.Store(usuarioBd.Apelido, &canal)
        return &canal, true
    }
}
// Buscando um usuario
func (g *GerenciadorCanaisSSE) buscarCanal(nomeCliente string) (*CanalSSE, bool) {
    if canal, existe := g.canais.Load(nomeCliente); existe {
        return canal.(*CanalSSE), true
    }
    return nil, false
}
//removendo um usuario
func (g *GerenciadorCanaisSSE) removerCanal(id string) error {
    if canal, existe := g.buscarCanal(id); existe {
        close(canal.Canal)
        g.canais.Delete(id)
        return nil
    }
    return nil
}
type SalaBD struct {
    Id int64 `json:"id"`
    Nome string `json:"nome"`
    Ativo bool `json:"ativo"`
}
type GerenciadorSalaBd struct {
    banco *database.Db
}
func newGerenciadorSalaBd() *GerenciadorSalaBd {
    return &GerenciadorSalaBd{
        banco: database.ConnectionConstructor(),
    }
}
func (sdb *GerenciadorSalaBd) criarSala(saladest *SalaBD) error {
    insertQuery := fmt.Sprintf("INSERT INTO sala (nome) VALUES ('%s')", saladest.Nome)
    var result sql.Result
    var err error
    result, err = sdb.banco.ExecAndLog(insertQuery)
    if err != nil {
        return err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    saladest.Id = id
    return nil
}
func (sdb *GerenciadorSalaBd) removerSala(sala *SalaBD) error {
    updateQuery := fmt.Sprintf("UPDATE sala SET ativo = 0 WHERE id = %d", sala.Id)
    _, err := sdb.banco.ExecAndLog(updateQuery)
    if err != nil {
       return err
    }
    return nil
}
func (sdb *GerenciadorSalaBd) buscarSala(nomesala string ,saladest *SalaBD) error {
    query := fmt.Sprintf("SELECT id, nome, ativo FROM sala WHERE nome = '%s' LIMIT 1", nomesala)
    linha := sdb.banco.QueryRowAndLog(query)
    if linha.Err() != nil {
        return linha.Err()
    }
    if err := linha.Scan(&saladest.Id, &saladest.Nome, &saladest.Ativo); err != nil {
        return err
    }
    return nil
}
func (sdb *GerenciadorSalaBd) adicionarUsuarioSala(usuarioSala *UsuarioSalaBD) error {
    insertQuery := fmt.Sprintf("INSERT INTO usuariosala (idusuario, idsala) VALUES (%d, %d)", usuarioSala.IdUsuario, usuarioSala.IdSala)
    var err error
    var resultado sql.Result
    resultado, err = sdb.banco.ExecAndLog(insertQuery)
    if  err != nil {
        return err
    }
    id, err := resultado.LastInsertId()
    if err != nil {
        return err
    }
    usuarioSala.Id = id
    return nil
}
func (sdb *GerenciadorSalaBd) removerUsuarioSala(sala *SalaBD, usuario *UsuarioBD) error {
    updateQuery := fmt.Sprintf("UPDATE usuariosala SET ativo = 0 WHERE idsala = %d AND idusuario = %d", sala.Id, usuario.Id)
    _, err := sdb.banco.ExecAndLog(updateQuery)
    if err != nil {
       return err
    }
    return nil
}
func (sdb *GerenciadorSalaBd) buscarUsuarioSala(idsala int64, idusuario int64, destSalaUsusario *UsuarioSalaBD) error {
    query := fmt.Sprintf("SELECT id ,idsala, idusuario, ativo FROM usuariosala WHERE idsala = %d AND idusuario = %d", idsala, idusuario)
    linha := sdb.banco.QueryRowAndLog(query)
    if linha.Err() != nil {
        return linha.Err()
    }
    if err := linha.Scan(&destSalaUsusario.Id, &destSalaUsusario.IdSala, &destSalaUsusario.IdUsuario, &destSalaUsusario.Ativo); err != nil{
        return err
    }
    return nil
}
func (sdb *GerenciadorSalaBd) ativarUsuarioSala(idUsuarioSala int64) error {
    updateQuery := fmt.Sprintf("UPDATE usuariosala SET ativo = 1 WHERE id = %d", idUsuarioSala)
    _, err := sdb.banco.ExecAndLog(updateQuery)
    if err != nil {
       return err
    }
    return nil
}
func (sdb *GerenciadorSalaBd) ativarSala(idsala int64) error {
    updateQuery := fmt.Sprintf("UPDATE sala SET ativo = 1 WHERE id = %d", idsala)
    _, err := sdb.banco.ExecAndLog(updateQuery)
    if err != nil {
       return err
    }
    return nil}
func (sdb *GerenciadorSalaBd) buscarUsuario(apelidousuario string ,usuariodest *UsuarioBD) error {
    query := fmt.Sprintf("SELECT id, nome, apelido FROM usuario WHERE apelido = '%s'", apelidousuario)
    linha := sdb.banco.QueryRowAndLog(query)
    if linha.Err() != nil {
        return linha.Err()
    }
    if err := linha.Scan(&usuariodest.Id, &usuariodest.Nome, &usuariodest.Apelido); err != nil {
        return err
    }
    return nil
}
func (sdb *GerenciadorSalaBd) adicionarUsuario(usuario *UsuarioBD) error {
    insertQuery := fmt.Sprintf("INSERT INTO usuario (nome, apelido) VALUES ('%s', '%s')", usuario.Nome, usuario.Apelido)
    var err error
    var res sql.Result
    res, err = sdb.banco.ExecAndLog(insertQuery)
    if err != nil {
        return err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return err
    }
    usuario.Id = id
    return nil
}
func (sdb *GerenciadorSalaBd) buscarSalasDeUsuario(idUsuario int64) ([]SalaBD ,error) {
    var salas []SalaBD
    query := fmt.Sprintf("SELECT s.id, s.nome FROM sala s INNER JOIN usuariosala us on us.idsala = s.id WHERE us.idusuario = %d", idUsuario)
    linhas, err := sdb.banco.QueryAndLog(query)
    defer linhas.Close()
    if err != nil {
        if err == sql.ErrNoRows {
            return salas, nil
        }
        return salas, err
    }
    for linhas.Next() {
        var sala SalaBD
        if err := linhas.Scan(&sala.Id, &sala.Nome); err != nil {
            return salas, nil
        }
        salas = append(salas, sala)
    }
    return salas, nil
}
func (sdb *GerenciadorSalaBd) buscarUsuariosDeSala(idSala int64) ([]string ,error) {
    var apelidos []string
    query := fmt.Sprintf("SELECT u.apelido FROM usuario u INNER JOIN usuariosala us on us.idusuario = u.id INNER JOIN sala s on s.id = us.idsala WHERE s.id = %d AND us.ativo = 1", idSala)
    linhas, err := sdb.banco.QueryAndLog(query)
    defer linhas.Close()
    if err != nil {
        if err == sql.ErrNoRows {
            return apelidos, nil
        }
        return apelidos, err
    }
    for linhas.Next() {
        var apelido string
        if err := linhas.Scan(&apelido); err != nil {
            return apelidos, nil
        }
        apelidos = append(apelidos, apelido)
    }
    return apelidos, nil
}
type UsuarioSalaBD struct {
    Id int64
    IdUsuario int64
    IdSala int64
    Ativo bool
}
type Sala struct {
    NomeSala string `json:"nomesala"`
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
func (s *Sala) transmitir(info *InfoSSE, gerenciadorCanais *GerenciadorCanaisSSE) {
    for _, cliente := range s.ClientesSala {
        if canal, existe := gerenciadorCanais.buscarCanal(cliente); existe {
            canal.Canal <- info
        }
    }
}
type GerenciadorSalas struct {
   // Salas sync.Map `json:"salas"`
    Salas map[string]*Sala `json:"salas"`
}
func (gsc *GerenciadorSalas) criarSala(nomeSala string) *Sala {
    if _, existe := gsc.Salas[nomeSala]; existe {
        return gsc.Salas[nomeSala]
    }
    var clientes []string
    sala := Sala{
        NomeSala: nomeSala,
        ClientesSala: clientes,
    }
    gsc.Salas[nomeSala] = &sala
    return &sala
}
func (gsc *GerenciadorSalas) buscarSala(nomeSala string) (*Sala, bool) {
    if sala, existe := gsc.Salas[nomeSala]; existe {
        return sala, true
    }
    var sala Sala
    return &sala, false
}
func (gs *GerenciadorSalas) removerSala(nomeSala string) {
    if _, existe := gs.Salas[nomeSala]; existe {
        delete(gs.Salas, nomeSala)
    }
}
type SalaComUsuariosApi struct {
    Id int64 `json:"id"`
    Nome string `json:"nome"`
    Usuarios []string `json:"usuarios"`
}
type LoginInput struct {
    Apelido string `json:"apelido" form:"apelido" bind:"required"`
}
func HandlerSSE(router *gin.Engine) {
    gerenciadorCanaisSSE := GerenciadorCanaisSSE{}
    gerenciadorSalasSSE := GerenciadorSalas{
        Salas: make(map[string]*Sala),
    }
    gerenciadorDb :=  newGerenciadorSalaBd()
    router.GET("/chat/:apelidousuario/salas", func(c *gin.Context) {
        var usuario UsuarioBD
        if err := gerenciadorDb.buscarUsuario(c.Param("apelidousuario"), &usuario); err != nil {
            if err == sql.ErrNoRows {
                c.AbortWithStatus(404)
                return
            }
            c.AbortWithStatus(500)
            return
        }
        salas, err := gerenciadorDb.buscarSalasDeUsuario(usuario.Id)
        var respostaSalasApi []SalaComUsuariosApi
        if err != nil {
            if err == sql.ErrNoRows {
                c.JSON(200, gin.H{
                    "status":"sucesso",
                    "salas": respostaSalasApi,
                })
                return
            }
            c.AbortWithStatus(500)
            return
        }
        for _, sala := range salas {
            apelidos, err := gerenciadorDb.buscarUsuariosDeSala(sala.Id)
            if err != nil {
                c.AbortWithStatus(500)
                return
            }
            respostaSalasApi = append(respostaSalasApi, SalaComUsuariosApi{
                Nome: sala.Nome,
                Usuarios: apelidos,
            })
        }
        c.JSON(200, gin.H{
            "status":"sucesso",
            "salas": respostaSalasApi,
        })
        return
    })
    router.GET("/chat/sse/:apelidousuario/entrar/:nomesala", func(c *gin.Context) {
        usuarioBd := UsuarioBD{}
        if err := gerenciadorDb.buscarUsuario(c.Param("apelidousuario"), &usuarioBd); err != nil {
            c.JSON(404, gin.H{
                "status":"falha",
                "mensagem": "Usuario nao encontrado",
            })
            fmt.Println("Usuario nao encontrado")
            return
        }
        _, existe := gerenciadorCanaisSSE.buscarCanal(usuarioBd.Apelido)
        if !existe {
            fmt.Println("Cliente não tem conexao de sse aberta")
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"Cliente não tem conexao de sse aberta",
            })
            return
        }
        sala, existe := gerenciadorSalasSSE.buscarSala(c.Param(("nomesala")))
        if !existe {
            sala = gerenciadorSalasSSE.criarSala(c.Param("nomesala"))
        }
        salabd := SalaBD {}
        if  err := gerenciadorDb.buscarSala(c.Param("nomesala"), &salabd); err != nil {
            if err == sql.ErrNoRows {
                salabd.Nome = c.Param("nomesala")
                if err = gerenciadorDb.criarSala(&salabd); err != nil {
                    gerenciadorSalasSSE.removerSala(sala.NomeSala)
                    fmt.Println("Erro ao criar sala")
                    c.AbortWithStatus(500)
                    return
                }
            } else {
                fmt.Println("Erro ao buscar sala", err)
                gerenciadorSalasSSE.removerSala(sala.NomeSala)
                c.AbortWithStatus(500)
                return
            }
        }
        if !salabd.Ativo {
            if err := gerenciadorDb.ativarSala(salabd.Id); err != nil {
                fmt.Println("Erro ao ativar sala")
                c.AbortWithStatus(500)
                return
            }
        }
        if sala.estaEmSala(usuarioBd.Apelido) {
            c.AbortWithStatus(204)
            return
        }
        usuarioSala := UsuarioSalaBD{
            IdSala: salabd.Id,
            IdUsuario: usuarioBd.Id,
        }
        err := gerenciadorDb.buscarUsuarioSala(usuarioSala.IdSala, usuarioSala.IdUsuario, &usuarioSala)
        if err == sql.ErrNoRows {
            if err == sql.ErrNoRows {
                if err := gerenciadorDb.adicionarUsuarioSala(&usuarioSala); err != nil{
                    fmt.Println("Erro ao adiiconar usuario sala")
                    c.AbortWithStatus(500)
                    return
                }
            }
        }
        if !usuarioSala.Ativo {
            if err := gerenciadorDb.ativarUsuarioSala(usuarioSala.Id); err != nil {
                fmt.Println("Erro ao ativar usuarioSala")
                fmt.Println(err)
                c.AbortWithStatus(500)
                return
            }
        }
        sala.adicionarCliente(usuarioBd.Apelido)
        infoEntrouEmSala := newInfoSSE("entrou-chat", newConteudo(c.Param("apelidousuario") + "entrou em uma sala", usuarioBd.Apelido, sala.NomeSala))
        go sala.transmitir(infoEntrouEmSala, &gerenciadorCanaisSSE)
        c.JSON(200, gin.H{
            "status":"sucesso",
            "mensagem": "sala criado e adentrada :)",
        })
        return
    })
    router.POST("/chat/sse/:apelidousuario/sala/:nomesala/enviar", func(c *gin.Context) {
        fmt.Println("Usuario pediu pra enviar mensagem", c.Param("apelidousuario"), c.Param("nomesala"))
        _, existe := gerenciadorCanaisSSE.buscarCanal(c.Param("apelidousuario"))
        if !existe {
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"Cliente não tem conexao de sse aberta",
            })
            return
        }
        var salaBd SalaBD
        if err := gerenciadorDb.buscarSala(c.Param("nomesala"), &salaBd); err != nil {
            if err == sql.ErrNoRows {
                fmt.Println("Servidor nao encontrou a sala")
                c.AbortWithStatus(404)
                return
            }
            c.AbortWithStatus(500)
            return
        }
        if !salaBd.Ativo {
            fmt.Println("Sala nao esta ativa")
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"nao é permitido enviar mensagens a salas nao ativas",
            })
            return
        }
        var usuarioBd UsuarioBD
        if err := gerenciadorDb.buscarUsuario(c.Param("apelidousuario"), &usuarioBd); err != nil {
             if err == sql.ErrNoRows {
                fmt.Println("Servidor nao encontrou usuario no banco", c.Param("apelidousuario"))
                c.AbortWithStatus(404)
                return
            }
            c.AbortWithStatus(500)
            return
        }
        var usuarioSala UsuarioSalaBD
        if err := gerenciadorDb.buscarUsuarioSala(salaBd.Id, usuarioBd.Id, &usuarioSala); err != nil {
            if err == sql.ErrNoRows {
                c.JSON(400, gin.H{
                    "status":"falha",
                    "mensagem": "usuario não é da sala",
                })
                return
            }
            fmt.Println("Erro ao buscar usuariosala")
            c.AbortWithStatus(500)
            return
        }
        if !usuarioSala.Ativo {
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem": "usuario não é membro ativo da sala",
            })
            return
        }
        sala, existe := gerenciadorSalasSSE.Salas[salaBd.Nome]
        if !existe {
            sala = gerenciadorSalasSSE.criarSala(salaBd.Nome)
        }
        if !sala.estaEmSala(usuarioBd.Apelido) {
            sala.adicionarCliente(usuarioBd.Apelido)
        }
        if len(c.Query("msg")) < 1 {
            c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"Conteudo de mensagem deve conter um valor!",
            })
            return
        }
        infoMensagemEnviadaCanal := newInfoSSE("chat-nova-mensagem", newConteudo(c.Query("msg"), usuarioBd.Apelido, salaBd.Nome))
        go sala.transmitir(infoMensagemEnviadaCanal, &gerenciadorCanaisSSE)
        c.AbortWithStatus(201)
        return
    })
    router.GET("/sse/:apelidousuario", func(c *gin.Context) {
        //buscar usuario em banco
        fmt.Println(c.Param("apelidousuario"), "Pediu pra abrir conexao")
        canal, existe := gerenciadorCanaisSSE.buscarCanal(c.Param("apelidousuario"))
        if existe {
            info := newInfoSSE("shutdown-sse", newConteudo("Printar aviso para utilizar chat em apenas uma aba", c.Param("apelidousuario"),""))
            canal.Canal <- info
            // limite de um chat aberto por usuario
            /*c.JSON(400, gin.H{
                "status":"falha",
                "mensagem":"usuario já possui sse aberto",
            })*/
        }
        //busacar usuario banco
        banco := database.ConnectionConstructor()
        query := fmt.Sprintf("SELECT id, nome, apelido FROM usuario WHERE apelido = '%s' LIMIT 1", c.Param("apelidousuario"))
        linha := banco.QueryRowAndLog(query)
        var usuarioBd UsuarioBD
        if err := linha.Scan(&usuarioBd.Id, &usuarioBd.Nome, &usuarioBd.Apelido); err != nil {
            if err.Error() == sql.ErrNoRows.Error() {
                fmt.Println("Erro ao buscar usuario no banco", c.Param("apelidousuario"))
                c.AbortWithStatus(404)
                return
            }
            fmt.Println("Erro ao escanear linha do banco")
            c.AbortWithStatus(500)
            return
        }
        canal, ok := gerenciadorCanaisSSE.criarCanal(usuarioBd)
        if !ok {
            c.AbortWithStatus(500)
        }
        fmt.Println("Canal criado >> " + canal.UsuarioBd.Apelido)
        c.Writer.Header().Set("Content-Type", "text/event-stream")
        c.Writer.Header().Set("Cache-Control", "no-cache")
        c.Writer.Header().Set("Connection", "keep-alive")
        go func() {
            <-c.Request.Context().Done()
            fmt.Println("Fechando canal de ",usuarioBd.Apelido)
            gerenciadorCanaisSSE.removerCanal(usuarioBd.Apelido)
        }()
        c.Stream(func(w io.Writer) bool {
            canal, _ = gerenciadorCanaisSSE.criarCanal(usuarioBd)
            if msg, ok := <- canal.Canal; ok {
                canal.gerenciarEventos(msg, c)
                return true
            }
            return false
        })
        return
    })
    router.GET("/chat/usuario/login", func (c *gin.Context) {
        c.HTML(200, "login.tmpl", gin.H{})
        return
    })
    router.POST("/chat/usuario/login", func (c *gin.Context) {
        var loginInput LoginInput
        if err := c.ShouldBind(&loginInput); err != nil {
            c.AbortWithStatus(404)
            return
        }
        var usuarioBd UsuarioBD
        if err := gerenciadorDb.buscarUsuario(loginInput.Apelido, &usuarioBd); err != nil {
            if err == sql.ErrNoRows {
                c.AbortWithStatus(404)
                return
            }
            c.AbortWithStatus(500)
            return
        }
        c.Redirect(302, fmt.Sprintf("/chat/%s/view", usuarioBd.Apelido))
        return
    })
    router.POST("/chat/api/usuario/login", func (c *gin.Context) {
        var loginInput LoginInput
        if err := c.ShouldBind(&loginInput); err != nil {
            c.AbortWithStatus(404)
            return
        }
        var usuarioBd UsuarioBD
        if err := gerenciadorDb.buscarUsuario(loginInput.Apelido, &usuarioBd); err != nil {
            if err == sql.ErrNoRows {
                c.AbortWithStatus(404)
                return
            }
            fmt.Println(err)
            c.AbortWithStatus(500)
            return
        }
        c.JSON(200, gin.H{
            "status": "sucesso",
            "usuario": usuarioBd,
        })
        return
    })
    router.GET("/chat/usuario/cadastrar", func(c *gin.Context) {
        c.HTML(200, "cadastrar.tmpl", gin.H{})
        return
    })
    router.POST("/chat/usuario/cadastrar", func(c *gin.Context) {
        var usuarioNovo UsuarioBD
        if err := c.ShouldBind(&usuarioNovo); err != nil {
            c.Redirect(302, "/chat/usuario/cadastrar?info=erro")
            return
        }
        gerenciadorDb := newGerenciadorSalaBd()
        if err := gerenciadorDb.adicionarUsuario(&usuarioNovo); err != nil {
            fmt.Println(err)
            fmt.Println("Erro ao adicionar usuario", usuarioNovo)
            c.Redirect(302, "/chat/usuario/cadastrar?info=erro ao adicionar")
            return
        }
        c.Redirect(302, fmt.Sprintf("/chat/%s/view", usuarioNovo.Apelido))
        return
    })
    router.POST("/chat/api/usuario/cadastrar", func(c *gin.Context) {
        var usuarioNovo UsuarioBD
        if err := c.ShouldBind(&usuarioNovo); err != nil {
            c.AbortWithStatus(500)
            return
        }
        if err := gerenciadorDb.adicionarUsuario(&usuarioNovo); err != nil {
            fmt.Println(err)
            fmt.Println("Erro ao adicionar usuario", usuarioNovo)
            c.AbortWithStatus(500)
            return
        }
        c.JSON(201, gin.H{
            "status": "sucesso",
            "usuario": usuarioNovo,
        })
        return
    })
    router.GET("/chat/:apelidousuario/view", func(c *gin.Context) {
        var usuario UsuarioBD
        if err := gerenciadorDb.buscarUsuario(c.Param("apelidousuario"), &usuario); err != nil {
            //quero direcionar para a rota de cadastro
            c.Redirect(302 , "/chat/usuario/cadastrar")
            return
        }
        c.HTML(200, "eventos.tmpl", gin.H{
            "title":"Chat",
            "apelidousuario": c.Param("apelidousuario"),
        })
        return
    })
    router.GET("/chat/sala/:nomesala/usuarios", func(c *gin.Context) {
        sala, existe := gerenciadorSalasSSE.buscarSala(c.Param("nomesala"))
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
    router.GET("/chat/usuario/:apelidousuario/sala/:nomesala/sair", func(c *gin.Context) {
        var usuarioBd UsuarioBD
        if err := gerenciadorDb.buscarUsuario(c.Param("apelidousuario"), &usuarioBd); err != nil {
            fmt.Println("Usuario nao encontrado")
            c.AbortWithStatus(404)
            return
        }
        var salaBd SalaBD
        if err := gerenciadorDb.buscarSala(c.Param("nomesala"), &salaBd); err != nil {
            fmt.Println("Sala nao encontrado")
            c.AbortWithStatus(404)
            return
        }
        canal, canalExiste := gerenciadorCanaisSSE.buscarCanal(usuarioBd.Apelido)
        if !canalExiste {
            fmt.Println("nao achou canal")
            c.JSON(404, gin.H{
                "status":"falha",
                "message":"canal nao encontrado",
            })
            return
        }
        infoDeixouSala := newInfoSSE("deixou-sala", newConteudo("", usuarioBd.Apelido, salaBd.Nome))
        sala, existe := gerenciadorSalasSSE.buscarSala(salaBd.Nome)
        if !existe && canalExiste {
            canal.Canal <- infoDeixouSala
            fmt.Println("nao achou sala")
        }
        var usuarioSala UsuarioSalaBD
        if err := gerenciadorDb.buscarUsuarioSala(salaBd.Id, usuarioBd.Id, &usuarioSala); err != nil {
            if err == sql.ErrNoRows {
                fmt.Println(err)
                fmt.Println("Usuario nao esta na sala")
                c.JSON(400,gin.H{
                    "status":"falha",
                    "mensagem":"usuariosala nao encontrado no banco",
                })
                return
            }
            fmt.Println(err, "Erro ao buscar usuario sala")
            c.AbortWithStatus(500)
            return
        }
        if usuarioSala.Ativo {
            if err := gerenciadorDb.removerUsuarioSala(&salaBd, &usuarioBd); err != nil {
                fmt.Println("Erro ao desativar usuario na sala")
                c.AbortWithStatus(500)
                return
            }
        }
        usuarios, err := gerenciadorDb.buscarUsuariosDeSala(salaBd.Id);
        if err != nil {
            fmt.Println("Erro ao buscar usuario de sala")
            c.AbortWithStatus(500)
            return
        }
        if len(usuarios) == 0 {
            if err := gerenciadorDb.removerSala(&salaBd); err != nil {
                fmt.Println("Erro ao desativar sala")
                c.AbortWithStatus(500)
                return
            }
        }
        if canalExiste {
            if !sala.estaEmSala(canal.UsuarioBd.Apelido) {
                fmt.Println("Usuario nao esta na sala em cache")
                canal.Canal <- infoDeixouSala
                c.AbortWithStatus(200)
                return
            }
            sala.removerCliente(canal.UsuarioBd.Apelido)
            canal.Canal <- infoDeixouSala
            if len(sala.ClientesSala) == 0 {
                gerenciadorSalasSSE.removerSala(sala.NomeSala)
            } else {
                for _, nomecliente := range sala.ClientesSala {
                    canalCliente, existe := gerenciadorCanaisSSE.buscarCanal(nomecliente)
                    if existe {
                        canalCliente.Canal <- infoDeixouSala
                    }
                }
            }
        }
        c.AbortWithStatus(200)
        return
    })
}
