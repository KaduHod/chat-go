package services

import (
	"chat/source/utils"
	"encoding/json"
	"errors"
	"fmt"
)

const (
    TIPO_AUTENTICACAO_PADRAO="m.login.password"
    API_LOGIN="/_matrix/client/r0/login"
    PROVEDOR_PADRAO="https://matrix.org"
    SYNC_API="/_matrix/client/r0/sync"
    MENSAGENS_DE_SALA="/_matrix/client/r0/rooms/%s/messages"
)

/*
https://spec.matrix.org/legacy/client_server/r0.2.0.html#get-matrix-client-r0-sync DOCS
*/
type ClienteMatrix struct {
    NomeDeUsuario string
    IdDeUsuario string
    UrlDoProvedor string
    Senha string
    Autenticacao AutenticacaoMatrix
    Salas []string
}

type AutenticacaoMatrix struct {
    AccessToken string `json:"access_token"`
    IdUsuario string `json:"user_id"`
    HomeServer string `json:"home_server"`
    IdDoDspositivo string `json:"device_id"`
    RefreshToken string
}

func NewClienteMatrix(nomeUsuario string, urlProvedor string, senha string) ClienteMatrix {
    var cliente ClienteMatrix
    cliente.NomeDeUsuario = nomeUsuario
    cliente.UrlDoProvedor = urlProvedor
    cliente.Senha = senha
    return cliente
}

func (c *ClienteMatrix) SetIdDeUsuario() {
    c.IdDeUsuario = fmt.Sprintf("@%s:matrix.org", c.NomeDeUsuario)
}

func (c *ClienteMatrix) Login() error {
    dadosCorpo := map[string]string {
        "password" : c.Senha,
        "type" : TIPO_AUTENTICACAO_PADRAO,
        "user" : c.IdDeUsuario,
    }
    corpoRequest, err := json.Marshal(dadosCorpo)
    if err != nil {
        return err
    }
    corpoResposta, respostaHttp, err := utils.RestPostRequest(fmt.Sprintf("%s%s", c.UrlDoProvedor, API_LOGIN), corpoRequest)
    if err != nil {
        return err
    }
    if respostaHttp.StatusCode != 200 {
        return errors.New("Erro ao se autenticar no servidor")
    }
    if err := json.Unmarshal(corpoResposta, &c.Autenticacao); err != nil {
        return err
    }
    fmt.Println(c)
    return nil
}

func (c *ClienteMatrix) BuscarSalas() error {
    url := fmt.Sprintf("%s%s?access_token=%s", c.UrlDoProvedor, SYNC_API, c.Autenticacao.AccessToken)
    corpoResposta, respostaHttp, err := utils.RestGetRequest(url)
    if err != nil {
        return err
    }
    if respostaHttp.StatusCode != 200 {
        return errors.New("Erro ao sincronizar")
    }
    var res map[string]interface{}
    if err := json.Unmarshal(corpoResposta, &res); err != nil {
        return err
    }
    salas, ok := res["rooms"].(map[string]interface{})
    if !ok {
        return errors.New("Erro ao pegar salas de resposta da requisicao")
    }
    salasQueEstouLogado, ok := salas["join"].(map[string]interface{})
    if !ok {
        return errors.New("Erro ao pegar salas logadas de resposta da requisicao")
    }
    for idSala := range salasQueEstouLogado {
        c.Salas = utils.AdicionaValorUnico(c.Salas, idSala)
    }
    return nil
}
func (c *ClienteMatrix) BuscarMensagensDeSala(idSala string) error {
    //url := fmt.Sprintf(MENSAGENS_DE_SALA, "")
}
