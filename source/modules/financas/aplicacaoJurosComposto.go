package financas

import (
	"chat/source/utils"
	"encoding/json"
	"log"
	"time"
)
const (
	Mensal = iota
	Trimestral
	Semestral
	Anual
	Nulo
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
		case Nulo:
			resultado = "nulo"
	}
	return resultado
} 

func PegaTipoFrequencia(opcao string) int {
	switch opcao {
	case "mensal":
		return Mensal
	case "anual":
		return Anual
	case "semestral":
		return Semestral
	case "trimestral":
		return Trimestral
	}	
	return Mensal
}

type AporteAplicacaoFinanceira struct {
	ValorAporte float64 `json:"valorAporte"`
	ValorAumentoAporte float64 `json:"valorAumentoAporte"`
	FrequenciaAumentoAporteDesc string `json:"frequenciaAumentoAporteDesc"`
	FrequenciaAumentoAporte int `json:"frequenciaAumentoAporte"`
	FrequenciaAporte int `json:"frequenciaAporte"`
}

type AplicacaoFinanceira struct {
	ValorInicial float64 `json:"valorInicial"`
	QuantidadeDeMeses int `json:"quantidadeDeMeses"`
	Aporte AporteAplicacaoFinanceira `json:"aporte"`
	Taxa float64 `json:"taxa"`
	ResultadoComJurosComposto float64 `json:"resultadoComJurosComposto"`
	ResultadoSemJurosComposto float64 `json:"resultadoSemJurosComposto"`
	ResultadoDePeriodos []ResultadoAplicacaoPeriodo `json:"resultadoDePeriodos"`
}

type ResultadoAplicacaoPeriodo struct {
	Identificador int `json:"identificador"`
	Data string `json:"data"`
	ValorComAporteMaisJurosCompostoFormatado string `json:"valorComJurosCompostoFormatado"`
	ValorComAporteSemJurosCompostoFormatado  string `json:"valorSemJurosCompostoFormatado"`
	ValorizacaoPeriodoFormatado string `json:"valorizacaoPeriodoFormatado"`
	ValorAportadoFormatado string `json:"valorAportadoFormatado"`
	ValorizacaoPeriodo float64 `json:"valorizacaoPeriodo"`
	ValorComAporteMaisJurosComposto float64 `json:"valorComAporteMaisJurosComposto"`
	ValorComAporteSemJurosComposto float64 `json:"valorComAporteSemJurosComposto"`
	ValorAportado float64 `json:"valorAportado"`
}

func (r *ResultadoAplicacaoPeriodo) formataValores() {
	r.ValorComAporteSemJurosCompostoFormatado  = utils.FloatParaValorMonetario(r.ValorComAporteSemJurosComposto)
	r.ValorComAporteMaisJurosCompostoFormatado = utils.FloatParaValorMonetario(r.ValorComAporteMaisJurosComposto)
	r.ValorizacaoPeriodoFormatado = utils.FloatParaValorMonetario(r.ValorizacaoPeriodo)
	r.ValorAportadoFormatado = utils.FloatParaValorMonetario(r.ValorAportado)
}

func (a *AplicacaoFinanceira) CalcularRendimento(valor *float64, valorSemJurosComposto *float64, data time.Time) ResultadoAplicacaoPeriodo {
	var resultadoPeriodo ResultadoAplicacaoPeriodo
	resultadoPeriodo.ValorizacaoPeriodo = a.calculaValorizacaoPeriodo(*valor)
	*valor = *valor + resultadoPeriodo.ValorizacaoPeriodo + a.Aporte.ValorAporte
	*valorSemJurosComposto = *valorSemJurosComposto + a.Aporte.ValorAporte
	resultadoPeriodo.Data = data.Format("2006-01-02")
	resultadoPeriodo.ValorAportado = a.Aporte.ValorAporte
	resultadoPeriodo.ValorComAporteSemJurosComposto = *valorSemJurosComposto
	resultadoPeriodo.ValorComAporteMaisJurosComposto = *valor
	resultadoPeriodo.formataValores()
	return resultadoPeriodo
}

func (a *AplicacaoFinanceira) IniciarSimulacao() {
	contador := 1
	valor := a.ValorInicial
	valorSemJurosComposto := a.ValorInicial
	dataInicial := time.Now()
	dataAcumuladora := dataInicial
	for contador <= a.QuantidadeDeMeses {
		switch a.Aporte.FrequenciaAporte {
		case Mensal:
			resultadoPeriodo := a.CalcularRendimento(&valor, &valorSemJurosComposto, dataAcumuladora)
			dataAcumuladora = dataAcumuladora.AddDate(0,1,0)
			resultadoPeriodo.Identificador = contador
			a.ResultadoDePeriodos = append(a.ResultadoDePeriodos, resultadoPeriodo)
		case Trimestral:
			if contador != 1 || contador % 3 != 0 {
				continue
			} 
			resultadoPeriodo := a.CalcularRendimento(&valor, &valorSemJurosComposto, dataAcumuladora)
			dataAcumuladora = dataAcumuladora.AddDate(0,3,0)
			resultadoPeriodo.Identificador = contador		
			a.ResultadoDePeriodos = append(a.ResultadoDePeriodos, resultadoPeriodo)
		case Semestral:
			if contador != 1 || contador % 6 != 0 {
				continue
			}
			resultadoPeriodo := a.CalcularRendimento(&valor, &valorSemJurosComposto, dataAcumuladora)
			dataAcumuladora = dataAcumuladora.AddDate(0,6,0) 
			resultadoPeriodo.Identificador = contador
			a.ResultadoDePeriodos = append(a.ResultadoDePeriodos, resultadoPeriodo)
		case Anual:
			if contador != 1 || contador % 12 != 0 {
				continue
			}
			resultadoPeriodo := a.CalcularRendimento(&valor, &valorSemJurosComposto, dataAcumuladora)
			dataAcumuladora = dataAcumuladora.AddDate(1,0,0) 
			resultadoPeriodo.Identificador = contador
			a.ResultadoDePeriodos = append(a.ResultadoDePeriodos, resultadoPeriodo)
		}
		contador++
		switch a.Aporte.FrequenciaAumentoAporte {
		case Mensal:
			a.Aporte.ValorAporte += a.Aporte.ValorAumentoAporte
		case Trimestral:
			if contador == 1 || contador % 3 != 0 {
				continue
			}
			a.Aporte.ValorAporte += a.Aporte.ValorAumentoAporte
		case Semestral:
			if contador == 1 || contador % 6 != 0 {
				continue
			}
			a.Aporte.ValorAporte += a.Aporte.ValorAumentoAporte
		case Anual:
			if contador == 1 || contador % 12 != 0 {
				continue
			}
			a.Aporte.ValorAporte += a.Aporte.ValorAumentoAporte
		}
	}
	a.ResultadoComJurosComposto = valor
	a.ResultadoSemJurosComposto = valorSemJurosComposto
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
