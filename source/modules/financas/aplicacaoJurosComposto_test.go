package financas_test

import (
	"chat/source/modules/financas"
	"testing"
)

func TestPrincipal(t *testing.T) {
	var aplicacao financas.AplicacaoFinanceira
	var aporte financas.AporteAplicacaoFinanceira
	t.Log(aporte, aplicacao)
	aplicacao.QuantidadeDeMeses = 12
	aporte.ValorAporte = 300.00
	aporte.ValorAumentoAporte = 100.00
	aporte.FrequenciaAporte = financas.Mensal
	aporte.FrequenciaAumentoAporte = financas.Anual
	aplicacao.Aporte = aporte
	aplicacao.ValorInicial = 0.0
	aplicacao.Taxa = 0.0097
	aplicacao.CalcularRendimento()
	t.Log(aporte.FrequenciaAporte)
	t.Fail();	
/*
	t.Log("Aquiii 3")
	t.Log(aplicacao.PrintarJsonComResultado())*/
}
