package database

import (
	"chat/source/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const (
	DUPLICATE_ENTRY = 1062
	DB_LOG_FILE = "/db.log"
)

type Db struct {
	Conn *sql.DB
}

func ConnectionConstructor() *Db {
	config := mysql.Config{
		User: os.Getenv("MYSQL_USERNAME"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net: "tcp",
		Addr: fmt.Sprintf("%s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")),
		DBName: os.Getenv("MYSQL_DATABASE"),
	}

	db, err := sql.Open("mysql", config.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr);
	}
	var database Db
	database.Conn = db
    database.Conn.SetMaxOpenConns(100)
	return &database
}
func (db *Db) FecharConexao() {
    if err := db.Conn.Close(); err != nil {
        fmt.Println("Erro ao fechar conexao")
        fmt.Println(err)
        panic(err)
    }
}

func (db *Db) ErroDeRegistroDuplicado(err error) bool {
    if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
        return true
    }
    return false
}

func (db *Db) QueryRowAndLog(query string, args ...any) *sql.Row {
	utils.Logger(DB_LOG_FILE, query, "DB LOG", true)
	var result *sql.Row
	if len(args) > 0 {
		result = db.Conn.QueryRow(query, args)
	} else {
		result = db.Conn.QueryRow(query)
	}
	if result.Err() != nil {
		utils.Logger(DB_LOG_FILE, result.Err().Error(), "DB LOG[EROOR]", true)
	}
	return result
}

func (db *Db) QueryAndLog(query string, args ...any) (*sql.Rows, error) {
	utils.Logger(DB_LOG_FILE, query, "DB LOG", true)
	var rows *sql.Rows
	var err error
	if len(args) > 0 {
		rows, err = db.Conn.Query(query, args)
	} else {
		rows, err = db.Conn.Query(query)
	}
	if err != nil {
		utils.Logger(DB_LOG_FILE, err.Error(), "DB LOG[EROOR]", true)
	}
	return rows, err
}

func (db *Db) ExecAndLog(query string, args ...any) (sql.Result, error) {
	utils.Logger(DB_LOG_FILE, query, "DB LOG", true)
	var result sql.Result
	var err error
	if len(args) > 0 {
		result, err = db.Conn.Exec(query, args)
	} else {
		result, err = db.Conn.Exec(query)
	}
	if err != nil {
		utils.Logger(DB_LOG_FILE, err.Error(), "DB LOG[EROOR]", true)
	}
	return result, err
}

type DummyStruct struct {
	Id int
}

func Exists(table string, id int) (bool, error) {
	var dummy DummyStruct
	conn := ConnectionConstructor()
	row := conn.QueryRowAndLog(fmt.Sprintf("SELECT id FROM %s WHERE id = %d LIMIT 1", table, id))
	defer conn.Conn.Close()
	if err := row.Err(); err != nil {
		return false, err
	}
	if err := row.Scan(&dummy.Id); err != nil {
		if err == sql.ErrNoRows {
			return false , nil
		}
		return false, err
	}
	return true, nil
}
