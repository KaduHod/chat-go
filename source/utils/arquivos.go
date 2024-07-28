package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
)
func EscreverEmArquivo(caminhoArquivo string, conteudo string) error {
    if _, err := os.Stat(caminhoArquivo); err != nil {
        _, err := os.Create(caminhoArquivo)
        if err != nil {
            return err
        }
    }
    if err := os.WriteFile(caminhoArquivo, []byte(conteudo), fs.ModeAppend); err != nil {
        return err
    }
    return nil
}
func SobrescreverArquivo(caminhoArquivo string, conteudo string) error {
    // Tenta criar o arquivo (isso trunca o arquivo se ele já existir)
    arquivoJson, err := os.Create(caminhoArquivo)
    if err != nil {
        fmt.Println("Erro ao criar ou abrir arquivo:", err)
        return err
    }

    // Garante que o arquivo será fechado ao final da função
    defer arquivoJson.Close()

    // Escreve o conteúdo no arquivo
    if _, err := arquivoJson.WriteString(conteudo); err != nil {
        fmt.Println("Erro ao sobrescrever arquivo:", err)
        return err
    }

    return nil
}
func LerArquivoJson[T any](caminhoArquivo string, dest *T) (error) {
    if _, err := os.Stat(caminhoArquivo); err != nil {
        fmt.Println("Arquivo não existe")
        return err
    }
    arquivoJson, err := os.Open(caminhoArquivo)
    if err != nil {
        fmt.Println("Erro ao abrir arquivo")
        return err
    }
    defer arquivoJson.Close()
    bytesArquivo, err := io.ReadAll(arquivoJson)
    if err != nil {
        fmt.Println("Erro ao ler conteudo de arquivo")
    }
    if err := json.Unmarshal(bytesArquivo, dest); err != nil {
        fmt.Println("Erro ao converter bytes para struct dest")
        return err
    }
    return nil
}
