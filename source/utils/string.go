package utils

import "regexp"

func VerificaPadrao(regex string, valor string) bool {
    re := regexp.MustCompile(regex)
    valoresEncontrados := re.FindStringSubmatch(valor)
    if len(valoresEncontrados) > 0 {
        return true
    }
    return false
}
