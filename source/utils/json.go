package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)
func StructParaJson[T any](str T) string {
	json, err := json.Marshal(str)
    if err != nil {
        fmt.Println("ERRO JSON: ",err)
        return ""
    }
	return string(json)
}
func JsonParaStruct[T any](str string, dest *T) {
	if err := json.Unmarshal([]byte(str), dest); err != nil{
		log.Fatal("Erro ao converter struct apra json")
	}
}
func Pretty(bts []byte) string {
	var preetyDest bytes.Buffer
	_ = json.Indent(&preetyDest, bts, "", "    ")
	return string(preetyDest.Bytes())
}
