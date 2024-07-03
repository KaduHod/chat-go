package utils

import (
	"encoding/json"
	"log"
)
func StructParaJson[T any](str T) string {
	json, err := json.Marshal(str)
	if err != nil {
		log.Fatal("Erro ao converter struct apra json")
	}
	return string(json)
}
func JsonParaStruct[T any](str string, dest *T) {
	if err := json.Unmarshal([]byte(str), dest); err != nil{
		log.Fatal("Erro ao converter struct apra json")
	}
}
