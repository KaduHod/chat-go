const iduser = document.getElementById("idusuario").value;
let SALA_SELECIONADA = false;
const salaMenuLateralContainer = document.getElementById("chat-salas")
const chatMensagensContainer = document.getElementById("chat-aberto-mensagens")
const listaMensagens = [];
const listaUsuarios = [];
const listaSalas = [];
const listaSalasUsuarios = [];
const listaSalasUsuariosMensagens = [];
class Utils {
    static async entrarSala(e){
        const nomeSala = document.getElementById('entrar-sala-valor').value
        if(nomeSala.length > 0) {
            const res = await fetch("/chat/sse/"+Utils.getIdUsuarioLogado()+"/entrar/"+nomeSala)
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
    static getIdUsuarioLogado() {
        return document.getElementById("idusuario").value;
    }
    static adicionaSalaAMenuLateral(sala) {
        const elementoHtml = document.createElement("div")
        elementoHtml.setAttribute("id", sala.id)
        elementoHtml.classList.add("item-card")
        elementoHtml.classList.add("rounded")
        elementoHtml.innerText = sala.nome
        elementoHtml.addEventListener("click", Utils.selecionarSala)
        salaMenuLateralContainer.appendChild(elementoHtml)
    }
    static selecionarSala(e) {
        const room = Sala.buscaPorId(e.target.id)
        if(!room) {
            throw new Error("Não foi possivel encontrar a sala")
        }
        SALA_SELECIONADA = room.id
        document.getElementById("titulo-janela-chat-aberta").innerText = room.nome
        const listaSum = SalaUsuarioMensagem.buscarPorSala(room.id)
        let idsMensagens = listaSum.map(v => v.idmensagem)
        const mensagens = listaMensagens.filter(msg => idsMensagens.includes(msg.id))
        const containerMensagens = document.getElementById("chat-aberto-mensagens")
        containerMensagens.innerHTML = ""
        mensagens.forEach(msg => containerMensagens.appendChild(msg.montaElemento()))
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
        return this.nome === Utils.getIdUsuarioLogado()
    }
}
class Mensagem {
    constructor(mensagem, idremetente){
        this.mensagem = mensagem
        this.id = "mensagem__"+Utils.generateUniqueId()
        this.data = Utils.formatarData();
        this.idremetente = idremetente
        this.alinhamento = Utils.getIdUsuarioLogado() === this.idremetente ? "self-start" : "self-end";
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
        divPrincipal.classList.add(this.alinhamento, 'flex', 'shadow-md', 'flex-col', 'rounded-md', 'p-2', 'max-w-fit', 'break-words', 'h-fit', 'bg-slate-200', 'my-1');

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
}
class Sala {
    constructor(nomeSala){
        this.nome = nomeSala
        this.id = `sala__${this.nome}`
    }
    static busca(nome) {
        return listaSalas.find(s => s.nome === nome)
    }
    static buscaPorId(id) {
        return listaSalas.find(s => s.id === id)
    }
    /*montaHtml(){
        this.html = `<div
            id="${this.id}"
            class="item-card rounded"
        >${this.nome}</div>`
        return this.html
    }*/
}
class SalaUsuario {
    constructor(iduser, idsala) {
        this.iduser = iduser;
        this.idsala = idsala;
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
try {
    const eventoSSE = new EventSource(`/sse/${iduser}`);
    eventoSSE.onerror = function(event) {
        console.error("Erro no SSE: ", event);
    };
    eventoSSE.addEventListener('entrou-chat', e => {
        const {sala, remetente} = JSON.parse(e.data).conteudo
        let room = Sala.busca(sala)
        if (!room) {
            room = new Sala(sala)
            listaSalas.push(room)
        }
        let usuario = listaUsuarios.find(user => user.nome === remetente)
        if(!usuario) {
            usuario = new Usuario(remetente)
            listaUsuarios.push(usuario)
        }
        let salaUsuario = listaSalasUsuarios.find(su => su.iduser === usuario.id && su.idsala === room.id)
        if(!salaUsuario) {
            salaUsuario = new SalaUsuario(usuario.id, room.id)
            listaSalasUsuarios.push(salaUsuario)
        }
        Utils.adicionaSalaAMenuLateral(room)
    })
    eventoSSE.addEventListener('chat-nova-mensagem', e => {
        e.data = JSON.parse(e.data)
        let {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
        const user = listaUsuarios.find(user => user.nome === remetente)
        if(!user) {
            throw new Error("Usuario não entrou na sala")
        }
        const room = Sala.busca(sala)
        if(!room) {
            throw new Error("Sala não encontrada")
        }
        const message = new Mensagem(mensagem, user.id)
        listaMensagens.push(message)
        const salaUsuarioMensagem = new SalaUsuarioMensagem(room.id, user.id, message.id)
        listaSalasUsuariosMensagens.push(salaUsuarioMensagem)
        console.log(SALA_SELECIONADA, room.nome)
        if(SALA_SELECIONADA === room.nome) {
            Utils.adicionaMensagemAoChat(message)
        }
    })
    eventoSSE.addEventListener('deixou-sala', e => {
        let {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
        return
    })
    eventoSSE.addEventListener('usuario-offline', e => {
        let {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
        return
    })
} catch (error) {
    console.log(error)
}
