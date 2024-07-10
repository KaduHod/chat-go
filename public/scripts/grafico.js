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
const taxa = document.getElementById('taxa')
const AMBIENTE = "PROD"
const URLBASE = AMBIENTE == "DEV" ? "http://localhost:3000" : "https://132f-2804-14c-87c4-82fb-83a5-9e0c-560d-fee5.ngrok-free.app" ;
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
	url = `${URLBASE}/financas/juroscomposto/simular?quantidadeDeMeses=${quantidadeAnos.value * 12}&valorInicial=${valorInicial.value}&taxa=${taxa.value}&valorAporte=${valorAporte.value}&frequenciaAporte=${frequenciaAporte.value}&frequenciaAumentoAporte=${frequenciaAumentoAporte.value}&valorAumentoAporte=${valorAumentoAporte.value}` 	
	console.log({url})
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
	const datas = resultadoDePeriodos.map((d) => d.data)
	const dataSets = [
		{
			label: 'Rendimentos com juros composto',
			data: valoresComJurosComposto, // Dados do eixo Y
			borderColor: 'rgba(11, 13, 82, 1)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		},
		{
			label: 'Rendimentos sem juros composto',
			data: valoresSemJurosComposto, // Dados do eixo Y
			borderColor: 'rgba(71, 11, 11, 1)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		},
		{
			label: 'Valorização do retorno da taxa de juros',
			data: valoresValorizacao, // Dados do eixo Y
			borderColor: 'rgba(10, 80, 10, 23)', // Cor da linha
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
	resultadoDiferenca.innerText = (ultimoRegistroDePeriodo.valorComAporteMaisJurosComposto - ultimoRegistroDePeriodo.valorComAporteSemJurosComposto).toLocaleString('pt-BR', {currency:"BRL", style:'currency'});
}
