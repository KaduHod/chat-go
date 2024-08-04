const iduser = document.getElementById("apelidousuario").value;
const salaMenuLateralContainer = document.getElementById("chat-salas")
const chatMensagensContainer = document.getElementById("chat-aberto-mensagens")
const botaoEnviarMensagem = document.getElementById("enviar-msg-botao")
const inputMensagem = document.getElementById("chat-text-input")
let listaSalasUsuariosMensagens = [];
let listaSalasUsuarios = [];
let listaMensagens = [];
let listaUsuarios = [];
let listaSalas = [];
let SALA_SELECIONADA = false;
class Utils {
    static async entrarSala(e){
        const nomeSala = document.getElementById('entrar-sala-valor').value
        if(nomeSala.length > 0) {
            const res = await fetch("/chat/sse/"+Usuario.getNomeUsuarioLogado()+"/entrar/"+nomeSala)
            if(res.status > 299) return
            SALA_SELECIONADA = nomeSala
            document.getElementById("titulo-janela-chat-aberta").innerText = nomeSala
        }
    }
    static getRandomHexColor() {
        const letters = '0123456789ABCDEF';
        let color = '#';
        for (let i = 0; i < 6; i++) {
            color += letters[Math.floor(Math.random() * 16)];
        }
        return color;
    }
    static generateUniqueId() {
        const randomPart = Math.random().toString(36).substring(2, 15);
        const timePart = Date.now().toString(36);
        return randomPart + timePart;
    }
    static formatarData() {
        const now = new Date();
        const day = String(now.getDate()).padStart(2, '0');
        const month = String(now.getMonth() + 1).padStart(2, '0'); // Janeiro é 0!
        const hours = String(now.getHours()).padStart(2, '0');
        const minutes = String(now.getMinutes()).padStart(2, '0');
        return `${day}/${month} ${hours}:${minutes}`;
    }
    static selecionarSala(e) {
        const room = Sala.buscaPorId(e.target.dataset.idsala)
        if(!room) {
            throw new Error("Não foi possivel encontrar a sala")
        }
        SALA_SELECIONADA = room.nome
        document.getElementById("titulo-janela-chat-aberta").innerText = room.nome
        const listaSum = SalaUsuarioMensagem.buscarPorSala(room.id)
        let idsMensagens = listaSum.map(v => v.idmensagem)
        const mensagens = listaMensagens.filter(msg => idsMensagens.includes(msg.id))
        const containerMensagens = document.getElementById("chat-aberto-mensagens")
        containerMensagens.innerHTML = ""
        mensagens.forEach(msg => containerMensagens.appendChild(msg.montaElemento()))
        botaoEnviarMensagem.addEventListener("click", Mensagem.enviarMensagem);
    }
    static removerMensagensChat() {
        const containerMensagens = document.getElementById("chat-aberto-mensagens")
        containerMensagens.innerHTML = ""
    }
    static removerNomeSalaTitulo() {
        document.getElementById("titulo-janela-chat-aberta").innerText = "Seleciona uma sala para visualizar as mensagens!"
    }
    static adicionaMensagemAoChat(mensagem) {
        const containerMensagens = document.getElementById("chat-aberto-mensagens")
        containerMensagens.appendChild(mensagem.montaElemento())
    }
}
document.getElementById("entrar-sala-botao").addEventListener("click", Utils.entrarSala)
class Usuario {
    constructor(nome, id = null, cor = null) {
        this.id = id ?? "usuario__"+Utils.generateUniqueId()
        this.nome = nome
        this.cor = cor ?? Utils.getRandomHexColor()
    }
    ehUsuarioLogado() {
        return this.nome === Usuario.getNomeUsuarioLogado()
    }
    static getNomeUsuarioLogado() {
        return document.getElementById("apelidousuario").value;
    }
    static busca(nome) {
        return listaUsuarios.find(u => u.nome == nome)
    }
}
class Mensagem {
    constructor(mensagem, usuario){
        this.mensagem = mensagem
        this.id = "mensagem__"+Utils.generateUniqueId()
        this.data = Utils.formatarData();
        this.idremetente = usuario.id
        this.alinhamento = Usuario.getNomeUsuarioLogado() !== usuario.nome ? "self-start" : "self-end";
    }
    static busca(id) {
        return listaMensagens.find(m => m.id === id)
    }
    /*montaHtml() {
        return `<div id="${this.id}" class="${this.alinhamento} flex shadow-md flex-col rounded-md p-2 max-w-fit break-words h-fit bg-slate-200 my-1">
                <div class="message-content">
                    ${this.mensagem}
                </div>
                <div class="flex justify-end mb-2">
                    <span class="mr-1 font-bold" style="color:${this.cor}">${this.remetente}</span>
                    <span class="text-slate-600">${Utils.formatarData()}</span>
                </div>
            </div>`
    }*/
    buscaDono() {
        return listaUsuarios.find(u => u.id === this.idremetente)
    }
    montaElemento() {
        const divPrincipal = document.createElement('div');
        divPrincipal.setAttribute('id', this.id);
        divPrincipal.classList.add(this.alinhamento, 'flex', 'shadow-md', 'flex-col', 'rounded-md', 'p-2', 'max-w-fit', 'break-words', 'h-fit', 'bg-slate-200', 'my-1', "animate-fade-in-move");

        // Cria o elemento <div> para a mensagem
        const messageContentDiv = document.createElement('div');
        messageContentDiv.classList.add('message-content');
        messageContentDiv.textContent = this.mensagem;

        // Adiciona o <div> da mensagem ao <div> principal
        divPrincipal.appendChild(messageContentDiv);

        // Cria o elemento <div> para o rodapé
        const footerDiv = document.createElement('div');
        footerDiv.classList.add('flex', 'justify-end', 'mb-2');

        // Cria o elemento <span> para o remetente
        const dono = this.buscaDono()
        const remetenteSpan = document.createElement('span');
        remetenteSpan.classList.add('mr-1', 'font-bold');
        remetenteSpan.setAttribute('style', `color:${dono.cor}`);
        remetenteSpan.textContent = dono.nome;

        // Cria o elemento <span> para a data
        const dataSpan = document.createElement('span');
        dataSpan.classList.add('text-slate-600');
        dataSpan.textContent = Utils.formatarData();

        // Adiciona os <span>s ao <div> do rodapé
        footerDiv.appendChild(remetenteSpan);
        footerDiv.appendChild(dataSpan);

        // Adiciona o <div> do rodapé ao <div> principal
        divPrincipal.appendChild(footerDiv);
        return divPrincipal;
    }
    static async enviarMensagem(e){
        const sala = Sala.buscaSalaSelecionada()
        if(!sala) {
            throw new Error("Sala nao encontrada!", {SALA_SELECIONADA})
        }
        const usuario = Usuario.busca(Usuario.getNomeUsuarioLogado())
        const conteudoMensagem = inputMensagem.value
        const resposta = await fetch(`/chat/sse/${usuario.nome}/sala/${sala.nome}/enviar?msg=${conteudoMensagem}`, {method:"POST"})
        if (resposta.status != 201) {
            const erroRequest = await resposta.json();
            throw new Error({"mensagem":"Erro ao enviar mensagem", err: erroRequest})
        }
    }
}
class Sala {
    constructor(nomeSala){
        this.nome = nomeSala
        this.id = `sala__${this.nome}`
    }
    montaElemento() {
        const elementoHtml = document.createElement("div")
        elementoHtml.setAttribute("id", `wrapper__${this.id}`)
        elementoHtml.classList.add("item-card", "rounded", "sala-menu-lateral", "animate-fade-in-move")
        const elementoNomeSala = document.createElement("div")
        elementoNomeSala.dataset.idsala = this.id
        elementoNomeSala.addEventListener("click", Utils.selecionarSala)
        elementoNomeSala.innerText = this.nome
        elementoNomeSala.classList.add("w-3/4", "flex", "items-center", "justify-center")
        const elementoSairSala = document.createElement("div")
        elementoSairSala.classList.add("flex", "w-1/4", "justify-center", "items-center")
        const iconeSair = document.createElement("img")
        iconeSair.setAttribute("src", "/public/sair.svg")
        iconeSair.classList.add("cursor-pointer", "w-[15px]")
        iconeSair.setAttribute("id", `sala__sair__${this.nome}`)
        iconeSair.dataset.nomesala = this.nome
        iconeSair.addEventListener("click", this.sair.bind(this))
        elementoHtml.appendChild(elementoNomeSala)
        elementoHtml.appendChild(iconeSair)
        return elementoHtml;
    }
    static busca(nome) {
        return listaSalas.find(s => s.nome === nome)
    }
    static buscaPorId(id) {
        return listaSalas.find(s => s.id === id)
    }
    static removerSalaDaLista(sala) {
        listaSalas = listaSalas.filter(s => s.nome !== sala.nome)
    }
    static salasMenuLateral() {
        return [...document.getElementsByClassName("sala-menu-lateral")].map(div => div.id)
    }
    static adicionaSalaAMenuLateral(sala) {
        if(document.getElementById(sala.id)) {
            throw new Error("Sala já esta no menu lateral!")
        }
        const elementoHtml = sala.montaElemento()
        salaMenuLateralContainer.appendChild(elementoHtml)
        document.getElementById("entrar-sala-valor").value = ""
    }
    static buscaSalaSelecionada() {
        return listaSalas.find(s => s.nome === SALA_SELECIONADA)
    }
    async buscarUsuariosServidor() {
        const resposta = await fetch(`/chat/sala/${this.nome}/usuarios`)
        if(resposta.status == 404) {
            throw new Error("Servidor nao encontrou a sala!")
        }
        return [...(await resposta.json()).clientes]
    }
    async sair() {
        const resposta = await fetch(`/chat/usuario/${Usuario.getNomeUsuarioLogado()}/sala/${this.nome}/sair`);
        if (resposta.status != 200) {
            const dado = await resposta.json()
            throw new Error("Erro ao enviar requisicao de sair de sala", {dado})
        }
    }
    static removerSalaMenuLateral(sala) {
        salaMenuLateralContainer.removeChild(document.getElementById(`wrapper__${sala.id}`))
        if(SALA_SELECIONADA == sala.nome) {
            Utils.removerMensagensChat()
            Utils.removerNomeSalaTitulo()
        }
    }
    /*montaHtml(){
        *
        this.html = `<div id="sala__1" class="item-card rounded flex">
            <div class="w-3/4 flex items-center justify-center">
            ${this.nome}
            </div>
            <div class="flex w-1/4 justify-center items-center">
            <img src="/public/sair.svg" class="cursor-pointer">
            </div>
            </div>`
        return this.html
    }*/
}
class SalaUsuario {
    constructor(iduser, idsala) {
        this.id = `salaUsuario__${Utils.generateUniqueId()}`
        this.iduser = iduser;
        this.idsala = idsala;
    }
    static busca({id ,iduser, idsala}) {
        const funcaoFiltro = !!iduser && !!idsala
            ? (item) => item.iduser == iduser && item.idsala == idsala
            : (!!iduser
                ? (item) => item.iduser == iduser
                : (item) => item.idsala == idsala
            )
        return listaSalasUsuarios.find(funcaoFiltro)
    }
    static removerPorId(id) {
        listaSalasUsuarios = listaSalasUsuarios.filter(s => s.id != id)
    }
}
class SalaUsuarioMensagem {
    constructor(idsala, idusuario, idmensagem) {
        this.idsala = idsala;
        this.idusuario = idusuario;
        this.idmensagem = idmensagem;
    }
    static buscarPorSala(idsala) {
        return listaSalasUsuariosMensagens.filter(sum => sum.idsala === idsala)
    }
}

const main = async () => {
    const respostaSalasUsuario = await fetch(`/chat/${iduser}/salas`)
    if (respostaSalasUsuario.status != 200) {
        throw new Error("Naõ foi possivel buscar as salas do usuario!")
    }
    const usuarioLogado = new Usuario(iduser)
    listaUsuarios.push(usuarioLogado)
    const dado = await respostaSalasUsuario.json()
    console.log(dado)
    const salas = [...dado.salas]
    salas.forEach(sala => {
        const room = new Sala(sala.nome)
        listaSalas.push(room)
        const salaUsuarioLogado = new SalaUsuario(usuarioLogado.id, room.id)
        listaSalasUsuarios.push(salaUsuarioLogado)
        if(!sala.usuarios) return;
        sala.usuarios.forEach(apelido => {
            if(apelido != usuarioLogado.nome) {
                const user = new Usuario(apelido)
                listaUsuarios.push(user)
                const salaUser = new SalaUsuario(user.id, room.id)
                listaSalasUsuarios.push(salaUser)
            }
        })
        Sala.adicionaSalaAMenuLateral(room)
    })
    const eventoSSE = new EventSource(`/sse/${iduser}`);
    eventoSSE.onerror = function(event) {
        console.error("Erro no SSE: ", event);
    };
    eventoSSE.addEventListener('entrou-chat', async e => {
        const {sala, remetente} = JSON.parse(e.data).conteudo
        let room = Sala.busca(sala)
        if (!room) {
            room = new Sala(sala)
            listaSalas.push(room)
        }
        let usuario = Usuario.busca(remetente)
        if(!usuario) {
            usuario = new Usuario(remetente)
            listaUsuarios.push(usuario)
        }
        let salaUsuario = SalaUsuario.busca({iduser: usuario.id, idsala: room.id})
        if(!salaUsuario) {
            salaUsuario = new SalaUsuario(usuario.id, room.id)
            listaSalasUsuarios.push(salaUsuario)
        }
        if (remetente === Usuario.getNomeUsuarioLogado()) {
            console.log("Aqui")
            // Busca usuarios ja logados na sala
            const usuariosJaLogadosEmSala = await room.buscarUsuariosServidor()
            usuariosJaLogadosEmSala.forEach(nomeUsuario => {
                let usuario = Usuario.busca(nomeUsuario)
                if(!usuario) {
                    usuario = new Usuario(nomeUsuario)
                    listaUsuarios.push(usuario)
                }
                let salaUsuario = SalaUsuario.busca({iduser: usuario.id})
                if(!salaUsuario) {
                    listaSalasUsuarios.push(new SalaUsuario(usuario.id, sala.id))
                }
            })
            if(!Sala.salasMenuLateral().includes(`wrapper__${room.id}`)) {
                Sala.adicionaSalaAMenuLateral(room)
            }
            if(!SALA_SELECIONADA) {
                SALA_SELECIONADA = room.nome
                botaoEnviarMensagem.addEventListener("click", Mensagem.enviarMensagem)
            }
        }
    })
    eventoSSE.addEventListener('chat-nova-mensagem', e => {
        e.data = JSON.parse(e.data)
        let {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
        const user = Usuario.busca(remetente)
        if(!user) {
            throw new Error("Usuario não entrou na sala")
        }
        const room = Sala.busca(sala)
        if(!room) {
            throw new Error("Sala não encontrada")
        }
        const message = new Mensagem(mensagem, user)
        listaMensagens.push(message)
        const salaUsuarioMensagem = new SalaUsuarioMensagem(room.id, user.id, message.id)
        listaSalasUsuariosMensagens.push(salaUsuarioMensagem)
        if(SALA_SELECIONADA === room.nome) {
            Utils.adicionaMensagemAoChat(message)
        }
    })
    eventoSSE.addEventListener('deixou-sala', e => {
        let {sala, remetente} = JSON.parse(e.data).conteudo
        const user = Usuario.busca(remetente)
        if (!user) {
            throw new Error("Usuario nao encontrado!")
        }
        const room = Sala.busca(sala)
        if (!room) {
            throw new Error("Sala nao encontrada!")
        }
        const salaUsuario = SalaUsuario.busca({iduser: user.id, idsala: room.id})
        if(!salaUsuario) {
            throw new Error("Usuario nao esta na sala para ser removido")
        }
        SalaUsuario.removerPorId(salaUsuario.id)
        Sala.removerSalaMenuLateral(room)
        Sala.removerSalaDaLista(room)
        if(room.nome == SALA_SELECIONADA) {
            SALA_SELECIONADA = false
            botaoEnviarMensagem.removeEventListener("click", Mensagem.enviarMensagem);
        }
        return
    })
    eventoSSE.addEventListener('usuario-offline', e => {
        let {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
        return
    })
}
try {
    main()
} catch (error) {
    console.log(error)
}

