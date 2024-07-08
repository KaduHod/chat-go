package financas

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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
	ResultadoComValorizacao float64 `json:"resultadocomvalorizacao"`
	ResultadoSemValorizacao float64 `json:"resultadosemvalorizacao"`
	MontanteResultadaDeValorizacao float64 `json:"montanteresultadadevalorizacao"`
	ResultadoDePeriodos []ResultadoAplicacaoPeriodo `json:"resultadodeperiodos"`
}
type ResultadoAplicacaoPeriodo struct {
	Data string `json:"data"`
	ValorizacaoPeriodo string `json:"valorizacaoperiodo"`
	MontanteComValorizacao string `json:"valorcomvalorizacao"`
	MontanteSemValorizacao string `json:"valorsemvalorizacao"`
}
type Dummy struct {
	Montante float64 `json:"montante"`
	MontanteComValorizacao float64 `json:"montante-com-valorizacao"`
	Valorizacao float64 `json:"valorizacao"`
	MontanteComAporteEValorizacao float64 `json:"montante-com-v-e-aporte"`
	MontanteSemValorizacao float64 `json:"montante-sem-valorizacao"`
	MontanteSemValorizacaoEComAporte float64 `json:"montante-sem-v-e-com-aporte"`
}
func (a *AplicacaoFinanceira) CalcularRendimento() {
	contador := 1
	montante := a.ValorInicial
	montanteSemValorizacao := a.ValorInicial
	dataInicial := time.Now()
	a.MontanteResultadaDeValorizacao = 0.0
	for contador <= a.QuantidadeDeMeses {
		switch a.Aporte.FrequenciaAporte {
		case Mensal:
			valorizacaoDoPeriodo := a.calculaValorizacaoPeriodo(montante)
			montante = montante + valorizacaoDoPeriodo
			montante = montante + a.Aporte.ValorAporte
			montanteSemValorizacao = montanteSemValorizacao + a.Aporte.ValorAporte
			a.MontanteResultadaDeValorizacao = montante
			a.ResultadoSemValorizacao = montanteSemValorizacao
			var resultadoPeriodo ResultadoAplicacaoPeriodo
			resultadoPeriodo.Data = dataInicial.Format("2006-01-02")
			resultadoPeriodo.MontanteSemValorizacao = resultadoPeriodo.paraValorMonetario(montanteSemValorizacao)
			resultadoPeriodo.MontanteComValorizacao = resultadoPeriodo.paraValorMonetario(montante)
			resultadoPeriodo.ValorizacaoPeriodo = resultadoPeriodo.paraValorMonetario(valorizacaoDoPeriodo)
			a.ResultadoDePeriodos = append(a.ResultadoDePeriodos, resultadoPeriodo)
			dataInicial = dataInicial.AddDate(0,1,0)
		case Trimestral:
		case Semestral:
		case Anual:
		}
		contador++
	}
	a.ResultadoComValorizacao = montante
	
}
func (a *AplicacaoFinanceira) PrintarJsonComResultado() string {
	jsonbytes, err := json.Marshal(a)
	if err != nil {
		log.Println("Erro ao converter resultado para json")
	}
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
