package utils
func AdicionaValorUnico(slice []string, value string) []string {
    // Verifica se o valor já existe na fatia
    for _, v := range slice {
        if v == value {
            return slice // Valor já existe, retorna a fatia sem modificações
        }
    }
    // Adiciona o valor à fatia
    return append(slice, value)
}
