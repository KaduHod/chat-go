package utils

import (
	"fmt"
	"log"
	"os"
	"time"

)

const (
	LOGGER_DIR = "/var/log/chat"
	DDMMYYY_hhmmss = "02/01/2006 15:04:05"
	TIME_LOCATION = "America/Sao_Paulo"
	ERROR_LOG_FILE = "/errors.log"
	DEBUG_LOG_FILE = "/debug.log"
)

func GetDateString() string {
	data := time.Now()
	location, err := time.LoadLocation(TIME_LOCATION); 
	if err != nil {
		fmt.Println("Erro ao alterar localização de data para o brazil")
	}
	data_brazil := data.In(location)
	return fmt.Sprintf("[%v]", data_brazil.Format(DDMMYYY_hhmmss))
}

func Logger(path string, msg string, title string, use_default_date_config bool) {
	file_name := LOGGER_DIR + path
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erro ao ABRIR arquivo " + file_name)
		return
	}
	var content []byte
	if use_default_date_config {
		content = []byte(GetDateString() + fmt.Sprintf("[%v] >>", title) + " " + msg + "\n")
	} else {
		content = []byte(fmt.Sprintf("[%v] >>", title) + " " + msg + "\n")
	}
	if _, err := file.Write(content); err != nil {
		fmt.Println("Erro ao ESCREVER em arquivo " + file_name)
		return
	}
	if err := file.Close(); err != nil {
		fmt.Println("Erro ao FECHAR arquivo " + file_name + " " + err.Error())
		log.Fatal(err)
	}
}
func Loga(msg string) {
	log.Println(msg)
	Logger("debug.log", msg, "Log rapido", true)
}
