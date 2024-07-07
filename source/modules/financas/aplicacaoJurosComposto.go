package financas

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)
const (
	Mensal = iota
	Trimestral
	Semestral
	Anual
)
func pegaTipoFrequencia(opcao int) string {
	var resultado string
	switch opcao {
		case Mensal:
			resultado = "mensal"
		case Anual:
			resultado = "anual"
		case Semestral:
			resultado = "semestral"
		case Trimestral:
			resultado = "trimestral"
	}
	return resultado
} 
type AporteAplicacaoFinanceira struct {
	ValorAporte float64 `json:"valoraporte"`
	ValorAumentoAporte float64 `json:"valoraumentoaporte"`
	FrequenciaAumentoAporteDesc string `json:"frequenciaaporte"`
	FrequenciaAumentoAporte int
	FrequenciaAporte int 
}
type AplicacaoFinanceira struct {
	ValorInicial float64 `json:"valorinicial"`
	QuantidadeDeMeses int `json:"quantidadedeemmeses"`
	Aporte AporteAplicacaoFinanceira `json:"aporte"`
	Taxa float64 `json:"taxa"`
	Resultado float64 `json:"resultado"`
	ResultadoSemTaxa float64 `json:"resultadosemtaxa"`
	MontanteResultadaDeValorizacao float64 `json:"montanteresultadadevalorizacao"`
	ResultadoDePeriodos []ResultadoAplicacaoPeriodo `json:"resultadodeperiodos"`
}
type ResultadoAplicacaoPeriodo struct {
	Data string `json:"data"`
	ValorizacaoPeriodo string `json:"valorizacaoperiodo"`
	MontanteComValorizacao string `json:"valorcomvalorizacao"`
	MontanteSemValorizacao string `json:"valorsemvalorizacao"`
}
func (a *AplicacaoFinanceira) CalcularRendimento() {
	contador := 1
	montante := a.ValorInicial
	a.MontanteResultadaDeValorizacao = 0.0
	for contador <= a.QuantidadeDeMeses {
		fmt.Println(contador)
		switch a.Aporte.FrequenciaAporte {
		case Mensal:
			var resultadoPeriodo ResultadoAplicacaoPeriodo
			valorizacaoDoPeriodo := a.calculaValorizacaoPeriodo(montante)
			valorSemValorizacao := montante + a.Aporte.ValorAporte
			valorComValorizacao := valorSemValorizacao + valorizacaoDoPeriodo
			a.MontanteResultadaDeValorizacao += valorComValorizacao
			resultadoPeriodo.ValorizacaoPeriodo = resultadoPeriodo.paraValorMonetario(valorizacaoDoPeriodo)
			resultadoPeriodo.MontanteComValorizacao = resultadoPeriodo.paraValorMonetario(a.MontanteResultadaDeValorizacao)
			resultadoPeriodo.MontanteSemValorizacao = resultadoPeriodo.paraValorMonetario(valorSemValorizacao)
			a.ResultadoDePeriodos = append(a.ResultadoDePeriodos, resultadoPeriodo)
			fmt.Println("Valorizacao: ",resultadoPeriodo.ValorizacaoPeriodo)
			contador++
		case Trimestral:
		case Semestral:
		case Anual:
		}
		contador++
	}
}
func (a *AplicacaoFinanceira) PrintarJsonComResultado() string {
	jsonbytes, err := json.Marshal(a)
	if err != nil {
		log.Println("Erro ao converter resultado para json")
	}
	log.Println(string(jsonbytes))
	return string(jsonbytes)
}
func (a *AplicacaoFinanceira) calculaValorizacaoPeriodo(valorIncrementado float64) float64 {
	return valorIncrementado * a.Taxa
}
func (a *AplicacaoFinanceira) CalcularRendimentoMensal() {}
func (a *AplicacaoFinanceira) CalcularRendimentoAnual() {}
func (a *AplicacaoFinanceira) CalcularRendimentoSemestral() {}
func (a *AplicacaoFinanceira) CalcularRendimentoTrimestral() {}
func (r *ResultadoAplicacaoPeriodo) paraValorMonetario(valor float64) string {
	valueStr := fmt.Sprintf("%.2f", valor)
	valueStr = strings.Replace(valueStr, ".", ",", 1)
	valueStr = "R$ " + valueStr
	return valueStr
}
