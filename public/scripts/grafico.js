const ctx = document.getElementById('meu-grafico')
const init = async () => {
	const res = await fetch(`http://localhost:3000/financas/juroscomposto/simular${window.location.search}`)
	let body = await res.json()
	console.log(body)
	const { resultadoDePeriodos } = body.aplicacao
	const valoresComJurosComposto = resultadoDePeriodos.map((d) => d.valorComAporteMaisJurosComposto);
	const valoresSemJurosComposto = resultadoDePeriodos.map((d) => d.valorComAporteSemJurosComposto)
	const datas = resultadoDePeriodos.map((d) => d.data)
	const dataSets = [
		{
			label: 'Rendimentos com juros composto',
			data: valoresComJurosComposto, // Dados do eixo Y
			borderColor: 'rgba(11, 13, 71, 1)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		},
		{
			label: 'Rendimentos sem juros composto',
			data: valoresSemJurosComposto, // Dados do eixo Y
			borderColor: 'rgba(71, 11, 11, 1)', // Cor da linha
			tension: 0.1, // Tensão da linha (0 para linhas retas, 0.1 para suavizar)
			fill: false // Não preencher área abaixo da linha
		}
	]
	new Chart(ctx, {
		type: 'line', // Tipo de gráfico: 'line' para gráfico de linha
		data: {
			labels: datas, // Rótulos do eixo X
			datasets: dataSets
		},
		options: {
			scales: {
				y: {
					beginAtZero: true // Começar o eixo Y no zero
				}
			}
		}
	});
}
init();
