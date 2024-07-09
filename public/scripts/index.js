const AMBIENTE = "PROD"//"PROD"
const BASEURLSERVER = AMBIENTE == "DEV" ? "localhost:3000" : "132f-2804-14c-87c4-82fb-83a5-9e0c-560d-fee5.ngrok-free.app";
const WEBSOCKETURL = AMBIENTE == "DEV" ? "ws://" + BASEURLSERVER + "/ws" : "wss://" + BASEURLSERVER + "/ws"
const nomeusertitulo = document.getElementById('nomeuser')
const div_resultado = document.getElementById("chat")
const text = document.getElementById('input')
const botao = document.getElementById('botao')
const botaocadastrar = document.getElementById('botaocadastrar')
const inputnome = document.getElementById('nome');
const inputapelido = document.getElementById('apelido');
const botaoconectar = document.getElementById("conectar")
let usuario = null
let ws = null
const inputEscondido = document.getElementById('hiddenususario')
function sleep() {
    return new Promise(resolve => setTimeout(resolve, 1300));
}
const botaoAmaMais = document.getElementById('botaodequemamamais')
botaoAmaMais.onclick = async () => {
	console.log('Clicou')
	console.log(inputEscondido, msgs)
	if(inputEscondido.value == "kaduhodi") {
		for(let c=0; c<msgs.length;c++) {
			let msg_ = {
				remetente: usuario.apelido,
				conteudo: msgs[c]
			}
			console.log(msg_)
			ws.send(JSON.stringify(msg_))
			await sleep()
		}
	} else {
		alert("Essa funcção é limitada apenas para quem ama mais [Carlos]")
	}
}

botaoconectar.onclick = () => {
	conectarWS();
}

botaocadastrar.onclick = async (evt) => {
	const avisaErro = (id = 1) => alert(`[${id}] Nome e apelido devem ter no minimo 8 caracteres`)
	const apelido = inputapelido.value.trim();
	let nome = apelido + " user name"
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
	inputEscondido.value = apelido
	usuario = user
	esconderElemento("cadastro")
	esconderElemento("textotutorial")
	nomeusertitulo.innerText = usuario.apelido == "Arnilloy" ? "Bem vinda, AMOR DA MINHA VIDA" : `Bem vindo ${usuario.apelido}` 
	mostrarElemento(nomeusertitulo.id)
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
	let url = AMBIENTE == "DEV" ? `http://${BASEURLSERVER}` : `https://${BASEURLSERVER}`;
	const res = await fetch(url + "/logar", {
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
	let cor = msg.remetente == "Arnilloy" ? "#46295A" : "#800000";
	return `<div class="message">
					<div class="message-content">
						${msg.conteudo}
					</div>
					<div class="message-footer">
		<span class="sender-name" style="color:${cor}">${msg.remetente}</span>
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
let msgs = [
    "Eu te amooooo!",
    "Você é o amor da minha vida.",
    "Não vivo sem você!",
    "Cada dia ao seu lado é um presente.",
    "Meu coração bate mais forte por você.",
    "Você me completa.",
    "Meu amor por você só cresce a cada dia.",
    "Você é meu tudo.",
    "Estou sempre pensando em você.",
    "Você é a melhor parte da minha vida.",
    "Você me faz tão feliz.",
    "Não consigo imaginar minha vida sem você.",
    "Você é a razão do meu sorriso.",
    "Eu te amo mais do que palavras podem expressar.",
    "Você ilumina meus dias.",
    "Sou grato por ter você na minha vida.",
    "Você é meu porto seguro.",
    "Cada momento com você é inesquecível.",
    "Você é minha alma gêmea.",
    "Minha vida é melhor com você.",
    "Você é minha inspiração.",
    "Você é minha melhor amiga e meu amor.",
    "Estou apaixonado por você mais a cada dia.",
    "Você é meu sonho realizado.",
    "Você é minha felicidade.",
    "Te amo com todo meu coração.",
    "Você é meu mundo.",
    "Você é minha razão de viver.",
    "Não posso viver sem você.",
    "Você é tudo para mim.",
    "Meu amor por você é eterno.",
    "Você é minha vida.",
    "Eu sou seu e você é minha.",
    "Você é a luz da minha vida.",
    "Você me faz sentir especial.",
    "Te amar é o melhor sentimento do mundo.",
    "Eu faria qualquer coisa por você.",
    "Você é a pessoa mais importante da minha vida.",
    "Meu amor por você é incondicional.",
    "Você é a melhor coisa que já me aconteceu.",
    "Eu sou a pessoa mais sortuda do mundo por ter você.",
    "Você é meu amor eterno.",
    "Você me faz querer ser uma pessoa melhor.",
    "Você é meu tudo.",
    "Te amo mais do que qualquer coisa.",
    "Você é minha alegria.",
    "Minha vida só faz sentido com você.",
    "Você é meu amor verdadeiro.",
    "Eu nunca te deixarei.",
    "Você é meu tesouro.",
    "Você é meu anjo.",
    "Você é a melhor parte de mim.",
    "Eu te amo de todo o coração.",
    "Você é minha paixão.",
    "Eu nunca me canso de você.",
    "Você é minha alma gêmea.",
    "Você é minha razão de ser.",
    "Você é minha estrela brilhante.",
    "Você é meu raio de sol.",
    "Você é minha felicidade infinita.",
    "Te amar é a coisa mais natural do mundo.",
    "Você é minha outra metade.",
    "Eu sou completamente louco por você.",
    "Você é meu desejo realizado.",
    "Você é meu amor de conto de fadas.",
    "Você é meu paraíso.",
    "Você é meu tudo.",
    "Você é minha razão de sorrir.",
    "Eu te amo com todo o meu ser.",
    "Você é meu maior presente.",
    "Você é meu destino.",
    "Você é meu sonho realizado.",
    "Você é minha vida.",
    "Eu te amo mais do que tudo.",
    "Você é minha esperança.",
    "Você é meu coração.",
    "Eu sou tão feliz por ter você.",
    "Você é minha vida.",
    "Te amo incondicionalmente.",
    "Você é minha razão de viver.",
    "Eu nunca te deixarei.",
    "Você é minha motivação.",
    "Você é meu amor eterno.",
    "Você é meu universo.",
    "Você é minha vida.",
    "Te amar é um privilégio.",
    "Você é minha razão de ser feliz.",
    "Você é meu tudo.",
    "Eu te amo mais a cada dia.",
    "Você é minha melhor escolha.",
    "Você é meu destino.",
    "Eu te amo de todo o coração.",
    "Você é minha vida.",
    "Você é meu futuro.",
    "Você é meu anjo da guarda.",
    "Você é meu tesouro mais precioso.",
    "Te amo infinitamente.",
    "Você é minha maior alegria.",
    "Você é minha felicidade completa.",
    "Você é meu mundo.",
    "Você é minha alma gêmea.",
    "Te amo mais do que posso expressar."
]

