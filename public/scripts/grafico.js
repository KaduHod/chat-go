const frequenciaAumentoAporte = document.getElementById('frequenciaAumentoAporte')
const valorAumentoAporte = document.getElementById('valorAumentoAporte')
const frequenciaAporte = document.getElementById('frequenciaAporte')
const quantidadeAnos = document.getElementById('quantidadeAnos')
const botaoSimular = document.getElementById('financaBotao')
const valorInicial = document.getElementById('valorInicial')
const valorAporte = document.getElementById('valorAporte')
const resultadoComJuros = document.getElementById('resultado_com_juros');
const resultadoSemJuros = document.getElementById('resultado_sem_juros');
const resultadoDiferenca = document.getElementById('diferenca_juros');
const maiorValorizacao = document.getElementById('maior_valorizacao');
const maiorValorAportado = document.getElementById('maior_valor_aportado')
const paraValorMonetario = v => v.toLocaleString('pt-BR', {currency:"BRL", style:'currency'});
const init = async () => {
	let url = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/eth.json"
	const res = await fetch(url);
	const {eth} = await res.json()
	document.getElementById('precoeth').innerText = paraValorMonetario(eth.brl)
}
init();
const taxa = document.getElementById('taxa')
const AMBIENTE = "DEV"
const chart = new Chart(document.getElementById('meu-grafico'), {
	type: 'line', // Tipo de gráfico: 'line' para gráfico de linha
	data: {
		labels: [], // Rótulos do eixo X
		datasets: []
	},
	options: {
		scales: {
			y: {
				beginAtZero: true // Começar o eixo Y no zero
			}
		}
	}
});
const ctx = document.getElementById('meu-grafico')
botaoSimular.onclick = async () => {
    if(quantidadeAnos.value > 80) {
        alert("Ta achando que é matusalem e vai viver mais de 1000 anos seu desgraçado! Vai travar meu pc!!!!!")
        return
    }
	url = `/financas/juroscomposto/simular?quantidadeDeMeses=${quantidadeAnos.value * 12}&valorInicial=${valorInicial.value}&taxa=${taxa.value}&valorAporte=${valorAporte.value}&frequenciaAporte=${frequenciaAporte.value}&frequenciaAumentoAporte=${frequenciaAumentoAporte.value}&valorAumentoAporte=${valorAumentoAporte.value}`
	const res = await fetch(url);
	if(res.status != 200) {
		alert("Erro ao fazer simulacao");
		return
	}
	let body = await res.json()
	const { resultadoDePeriodos } = body.aplicacao
	const valoresComJurosComposto = resultadoDePeriodos.map((d) => d.valorComAporteMaisJurosComposto);
	const valoresSemJurosComposto = resultadoDePeriodos.map((d) => d.valorComAporteSemJurosComposto)
	const valoresValorizacao = resultadoDePeriodos.map((d) => d.valorizacaoPeriodo)
	const valorAportado = resultadoDePeriodos.map((d) => d.valorAportado)
	const datas = resultadoDePeriodos.map((d) => d.data)
	const dataSets = [
		{
			label: 'Rendimentos com juros composto',
			data: valoresComJurosComposto, // Dados do eixo Y
			borderColor: 'rgba(245, 68, 55, 1)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		},
		{
			label: 'Rendimentos sem juros composto',
			data: valoresSemJurosComposto, // Dados do eixo Y
			borderColor: 'rgba(114, 208, 245, 1)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		},
		{
			label: 'Valorização do retorno da taxa de juros',
			data: valoresValorizacao, // Dados do eixo Y
			borderColor: 'rgba(68, 245, 124, 23)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		},
		{
			label: 'Valor aportado',
			data: valorAportado, // Dados do eixo Y
			borderColor: 'rgba(127, 136, 124, 23)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		},
	]
	chart.data.datasets = dataSets
	chart.data.labels = datas
	chart.update()
	let ultimoRegistroDePeriodo = resultadoDePeriodos[resultadoDePeriodos.length-1];
	resultadoComJuros.innerText = ultimoRegistroDePeriodo.valorComJurosCompostoFormatado;
	resultadoSemJuros.innerText = ultimoRegistroDePeriodo.valorSemJurosCompostoFormatado;
	maiorValorizacao.innerText = ultimoRegistroDePeriodo.valorizacaoPeriodoFormatado;
	maiorValorAportado.innerText = paraValorMonetario(ultimoRegistroDePeriodo.valorAportado)
	resultadoDiferenca.innerText = paraValorMonetario(ultimoRegistroDePeriodo.valorComAporteMaisJurosComposto - ultimoRegistroDePeriodo.valorComAporteSemJurosComposto);
}
