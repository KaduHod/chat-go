const AMBIENTE = "DEV"//"PROD"
const BASEURLSERVER = AMBIENTE == "DEV" ? "localhost:3000" : "37a1-2804-14c-87c4-82fb-8e69-ca42-e2b0-69fa.ngrok-free.app"
const WEBSOCKETURL = BASEURLSERVER + "/ws"

const div_resultado = document.getElementById("chat")
const text = document.getElementById('input')
const botao = document.getElementById('botao')
const botaocadastrar = document.getElementById('botaocadastrar')
const inputnome = document.getElementById('nome');
const inputapelido = document.getElementById('apelido');
let usuario = null
console.log({botaocadastrar})

botaocadastrar.onclick = async (evt) => {
	const avisaErro = (id = 1) => alert(`[${id}] Nome e apelido devem ter no minimo 8 caracteres`)
	const nome = inputnome.value.trim()
	const apelido = inputapelido.value.trim();
	if(nome == "" || apelido == ""){
		return avisaErro()
	}
	if(nome.length < 8 || apelido.length != 8) {
		return avisaErro(2)
	}
	if(!(user = await logarRequest({nome, apelido}))) {
		console.log({user})
		return alert("Erro ao logar :(")	
	}
	esconderElemento("cadastro")
}

const esconderElemento = id => {
	const elemento = document.getElementById(id)
	if(elemento.classList.contains('escondido')) return
	return elemento.classList.add('escondido')
}

const logarRequest = async (usuario) => {
	const res = await fetch("http://" + BASEURLSERVER + "/logar", {
		headers: {'Content-type':'application/json'},
		method: "POST",
		body: JSON.stringify(usuario)
	})
	const resposta = await res.json()
	console.log({res, resposta})
	usuario = resposta.usuario
	if(res.status !== 200){
		return false
	}
	return resposta.usuario
}

const conectarWS = async () => {
	const ws = new WebSocket(WEBSOCKETURL)
	botao.onclick = () => {
		ws.send(text.value)
		text.value = ""
	}

	ws.onclose = (evt) => {
		console.log("Conexão fechada")
		escreveMensagem("Conexao fechada com o servidor de chat");
	}

	document.addEventListener("keydown", e => {
		if (e.key === "Enter") {
			botao.click()
		}
	})

	ws.onmessage = (evt) => {
		console.log({evt});
		console.log(evt.data)
		escreveMensagem("[anônimo]  " + evt.data)
	}
}

const escreveMensagem = msg => {
	div_resultado.innerHTML += `<br> [${formatarData()}] ${msg}`
}

function formatarData() {
    const now = new Date();
    const day = String(now.getDate()).padStart(2, '0');
    const month = String(now.getMonth() + 1).padStart(2, '0'); // Janeiro é 0!
    const year = now.getFullYear();
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    return `${day}/${month}/${year} ${hours}:${minutes}`;
}
