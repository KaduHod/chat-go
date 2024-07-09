package financas_test

import (
	"bytes"
	"chat/source/modules/financas"
	"encoding/json"
	"fmt"
	"testing"
)

func TestPrincipal(t *testing.T) {
	var aplicacao financas.AplicacaoFinanceira
	var aporte financas.AporteAplicacaoFinanceira
	aplicacao.QuantidadeDeMeses = 12*30
	aporte.ValorAporte = 414.00
	aporte.ValorAumentoAporte = 100.00
	aporte.FrequenciaAporte = financas.Mensal
	aporte.FrequenciaAumentoAporte = financas.Anual
	aplicacao.Aporte = aporte
	aplicacao.ValorInicial = 2940.0
	aplicacao.Taxa = 0.0097
	aplicacao.IniciarSimulacao()
	var json_preety bytes.Buffer
	_ = json.Indent(&json_preety, []byte(aplicacao.PrintarJsonComResultado()), "", "\t")
	fmt.Println(string(json_preety.Bytes()))
	t.Fail()
}
