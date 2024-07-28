package utils

import (
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
