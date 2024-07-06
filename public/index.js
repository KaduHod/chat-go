const ws = new WebSocket("wss://37a1-2804-14c-87c4-82fb-8e69-ca42-e2b0-69fa.ngrok-free.app/ws")
const div_resultado = document.getElementById("chat")
const text = document.getElementById('input')
const button = document.getElementById('botao')
const botaocadastrar = document.getElementById('botaocadastrar')
const inputnome = document.getElementById('nome');
const inputapelido = document.getElementById('apelido');
botaocadastrar.onclick = async (evt) => {
	const avisaErro = () => alert('Nome e apelido devem ter no minimo 8 caracteres')
	const nome = inputnome.value.trim()
	const apelido = inputapelido.value.trim();
	console.log({
		len: nome.length,
		len2: apelido.length
	})
	if(nome == "" || apelido == ""){
		return avisaErro()
	}
	if(nome.length < 8 || apelido.length != 8) {
		return avisaErro()
	}
	fetch
}


button.onclick = () => {
	ws.send(text.value)
	text.value = ""
}
ws.onclose = (evt) => {
	console.log("Conexão fechada")
	escreveMensagem("Conexao fechada com o servidor de chat");
}
document.addEventListener("keydown", e => {
	if (e.key === "Enter") {
		button.click()
	}
})
ws.onmessage = (evt) => {
	console.log({evt});
	console.log(evt.data)
	escreveMensagem("[anônimo]  " + evt.data)
}

const escreveMensagem = msg => {
	div_resultado.innerHTML += `<br> [${formatDate()}] ${msg}`
}

function formatDate() {
    const now = new Date();
    const day = String(now.getDate()).padStart(2, '0');
    const month = String(now.getMonth() + 1).padStart(2, '0'); // Janeiro é 0!
    const year = now.getFullYear();
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    return `${day}/${month}/${year} ${hours}:${minutes}`;
}