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
}
