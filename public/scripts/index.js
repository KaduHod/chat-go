const AMBIENTE = "DEV"//"PROD"
const BASEURLSERVER = AMBIENTE == "DEV" ? "localhost:3000" : "37a1-2804-14c-87c4-82fb-8e69-ca42-e2b0-69fa.ngrok-free.app"
const WEBSOCKETURL = AMBIENTE == "DEV" ? "ws://" + BASEURLSERVER + "/ws" : "wss://" + BASEURLSERVER + "/ws"

const div_resultado = document.getElementById("chat")
const text = document.getElementById('input')
const botao = document.getElementById('botao')
const botaocadastrar = document.getElementById('botaocadastrar')
const inputnome = document.getElementById('nome');
const inputapelido = document.getElementById('apelido');
const botaoconectar = document.getElementById("conectar")
let usuario = null
let ws = null

botaoconectar.onclick = () => {
	conectarWS();
}

botaocadastrar.onclick = async (evt) => {
	const avisaErro = (id = 1) => alert(`[${id}] Nome e apelido devem ter no minimo 8 caracteres`)
	let nome = inputnome.value.trim()
	const apelido = inputapelido.value.trim();
	nome = apelido + " user name"
	if(nome == "" || apelido == ""){
		return avisaErro()
	}
	if(nome.length < 8 || apelido.length != 8) {
		return avisaErro(2)
	}
	if(!(user = await loginServico({nome, apelido}))) {
		console.log({user})
		return alert("Erro ao logar :(")	
	}
	usuario = user
	esconderElemento("cadastro")
	mostrarElemento('chatcontainer');
	conectarWS()
}

const mostrarElemento = id => {
	const elemento = document.getElementById(id)
	if(elemento.classList.contains('escondido')) elemento.classList.remove('escondido') 
	return
}

const esconderElemento = id => {
	const elemento = document.getElementById(id)
	if(elemento.classList.contains('escondido')) return
	return elemento.classList.add('escondido')
}

const loginServico = async (usuario) => {
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

const montaMensagem =  (msg) => {
	return `<div class="message">
					<div class="message-content">
						${msg.conteudo}
					</div>
					<div class="message-footer">
						<span class="sender-name">${msg.remetente}</span>
						<span class="message-time">${formatarData()}</span>
					</div>
				</div>`
}
function scrollChatToBottom() {
    var chatDiv = document.getElementById("chat");
    chatDiv.scrollTop = chatDiv.scrollHeight;
}
const conectarWS = async () => {
	ws = new WebSocket(WEBSOCKETURL)
	botao.onclick = () => {
		let msg = {
			remetente: usuario.apelido,
			conteudo: text.value
		}
		console.log({msg})
		ws.send(JSON.stringify(msg))
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
		let msg = JSON.parse(evt.data)
		escreveMensagem(montaMensagem(msg))
		scrollChatToBottom()
	}
}

const escreveMensagem = msg => {
	div_resultado.innerHTML += `${msg}`
}

function formatarData() {
    const now = new Date();
    const day = String(now.getDate()).padStart(2, '0');
    const month = String(now.getMonth() + 1).padStart(2, '0'); // Janeiro é 0!
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    return `${day}/${month} ${hours}:${minutes}`;
}
