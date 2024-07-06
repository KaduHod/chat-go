package services

import (
	"chat/source/database"
	"chat/source/entidades"
	"chat/source/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)
func LogarUsuario(c *gin.Context) {
	var usr entidades.Usuario
	if !utils.RequestBody[entidades.Usuario](&usr, c) {
		return
	} 
	if erros, ok := usr.ValidarAtributos(); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"mensagem": "erro",
			"erros": erros,
		})
		return
	}
	db := database.ConnectionConstructor()
	if ok := usr.BuscarUsuarioBanco(); !ok {
		sql := fmt.Sprintf("INSERT INTO usuario (nome, apelido) VALUES ('%s','%s')", usr.Nome, usr.Apelido)
		db.ExecAndLog(sql)
		if !usr.BuscarUsuarioBanco(){
			defer db.Conn.Close()
			utils.Logger("/debug.log", "Erro ao inserir usuario ao banco", "LOGAR USUARIO", true)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	defer db.Conn.Close()
	c.JSON(http.StatusOK, gin.H{
		"mensagem": "ok",
		"usuario": usr,
		"chave_secreta": utils.HashSha512(fmt.Sprintf("%d", usr.Id)),
	})
	return 
}
