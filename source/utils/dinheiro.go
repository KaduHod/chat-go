package utils

import (
	"strings"

	"github.com/dustin/go-humanize"
)

func FloatParaValorMonetario(quantia float64) string {
	valueStr := "R$ " + humanize.Commaf(quantia)
	posPonto := strings.Index(valueStr, ".")
	if posPonto == -1 {
		valueStr = valueStr + ".00"
	} else if len(valueStr) > posPonto + 3 {
		valueStr = valueStr[:posPonto+3]	
	}
	swapped := strings.ReplaceAll(valueStr, ",", "X")
	swapped = strings.ReplaceAll(swapped, ".", ",")
	swapped = strings.ReplaceAll(swapped, "X", ".")
	return swapped
}
