package entidades

import (
	"chat/source/database"
	"chat/source/utils"
	"fmt"
)

type Usuario struct {
	Id int `json:"id"`
	Nome string `json:"nome"`
	Apelido string `json:"apelido"`
}

func (u *Usuario) ValidarAtributos() ([]string, bool) {
	var mensagensErro []string	
	sucesso := true
	if len(u.Nome) < 8 {
		sucesso = false
		mensagensErro = append(mensagensErro, "Nome deve ter mais do que 7 caracteres")
	}
	if (len(u.Apelido) != 8) {
		sucesso = false
		mensagensErro = append(mensagensErro, "Apelido deve possuir 8 caracteres")
	}
	return mensagensErro, sucesso
}

func (u *Usuario) BuscarUsuarioBanco() bool {
	db := database.ConnectionConstructor()
	row := db.QueryRowAndLog(fmt.Sprintf("SELECT id, nome, apelido FROM usuario WHERE apelido = '%s'", u.Apelido));
	if row.Err() != nil {
		utils.Logger("/debug.log", "Erro ao buscar usuario", "BUSCAR USUARIO", true)
		defer db.Conn.Close()
		return false
	}
	defer db.Conn.Close()
	if err := row.Scan(&u.Id, &u.Nome, &u.Apelido); err != nil {
		utils.Logger("/debug.log", "Erro ao escanear resultado de banco para struct", "BUSCAR USUARIO", true)
		defer db.Conn.Close()
		return false

	}
	return true
}
