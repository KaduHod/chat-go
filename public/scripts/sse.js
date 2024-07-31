const iduser = document.getElementById("idusuario").value;
const salaItemContainer = document.getElementById("chat-salas")
const chatContainer = document.getElementById("chat-aberto-mensagens")
const salaitemConstructor = nomeSala => `<div id="sala___${nomeSala}" class="item-card rounded">${nomeSala}</div>`

const mensagensSalas = {}
const usuarios = {}
class Usuario {
    constructor(nome, id = null, cor = null) {
        this.id = id ?? "usuario__"+Utils.generateUniqueId()
        this.nome = nome
        this.cor = cor ?? Utils.getRandomHexColor()
        this.__id = nome
    }
    ehUsuarioLogado() {
        return this.__id == iduser
    }
    salva(){
        if(!usuarios[this.__id]) {
            usuarios[this.__id] = this
        }
    }
}
class Mensagem {
    constructor(input){
        this.remetente = Utils.getUsuario(input.nomeUsuario)
        this.mensagem = input.mensagem
        this.idsala = input.idsala
        this.alinhamento = input.remetente ?? this.remetente == iduser ? "self-end" : "self-start"
        this.cor = this.remetente.cor
        this.id = input.id ?? "mensagem__"+Utils.generateUniqueId()
    }
    montaHtml() {
        return `<div id="${this.id}" class="${this.alinhamento} flex shadow-md flex-col rounded-md p-2 max-w-fit break-words h-fit bg-slate-200 my-1">
                <div class="message-content">
                    ${this.mensagem}
                </div>
                <div class="flex justify-end mb-2">
                    <span class="mr-1 font-bold" style="color:${this.cor}">${this.remetente}</span>
                    <span class="text-slate-600">${Utils.formatarData()}</span>
                </div>
            </div>`
    }
    salva() {
        if(mensagensSalas[this.idsala]) {
            mensagensSalas[this.idsala] = []
        }
        mensagensSalas[this.idsala].push(this)
    }
    getSala(){
        return new Sala(salas[this.idsala].nome)
    }
}
class Sala {
    constructor(nomeSala){
        this.nome = nomeSala
        this.id = this.getId();
    }
    getMensagens(){
        if(!mensagensSalas[this.id]) return [];
        return mensagensSalas[this.id].map( msg => new Mensagem(msg) )
    }
    salva(){
        if(!salas[this.id]){
            salas[this.id] = this
        }
    }
    salvaMensagem(mensagem) {
        if(!mensagensSalas[this.id]) {
            mensagensSalas[this.id] = []
        }
        mensagensSalas[this.nome].push(mensagem)
    }
    adicionaMensagemEmContainer(mensagem) {
        chatContainer.innerHTML += mensagem.montaHtml()
    }
    getId(){
        return "sala__"+this.nome
    }
    getElemento() {
        return document.getElementById(this.id)
    }
    montaHtml(){
        this.html = `<div
            id="${this.id}"
            class="item-card rounded"
        >${this.nome}</div>`
        return this.html
    }
    selecionaSala() {
        GerenciadorSala.containerTituloChatAberto().innerText = this.nome
    }
}
const salas = {
    "sala__1": new Sala("1"),
    "sala__2": new Sala("2"),
    "sala__3": new Sala("3"),
}
class GerenciadorSala {
    static getSala(id){
        return new Sala(salas[id].nome);
    }
    static adicionaSalaMenu(idsala){
        const sala = GerenciadorSala.getSala(idsala)
        GerenciadorSala.containerSalasMenu().innerHTML += sala.montaHtml();
        sala.getElemento().onclick = GerenciadorSala.selecionaSalaEvento
        console.log(sala, sala.getElemento())
        sala.getElemento().click()
    }
    static adicionaMensagemChat(idmensagem) {

    }
    static containerSalasMenu() {
        return document.getElementById("chat-salas")
    }
    static containerMensagens() {
        return document.getElementById("chat-aberto-mensagens")
    }
    static containerTituloChatAberto() {
        return document.getElementById("titulo-janela-chat-aberta")
    }
    static selecionaSalaEvento(e){
        GerenciadorSala
            .getSala(this.id)
            .selecionaSala()
    }
}
class Utils {
    static getUsuario(nome){
        if(usuarios[nome]) return new Usuario(usuarios[nome].nome, usuarios[nome].id, usuarios[nome].cor)
        return new Usuario(nome)
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
    formatarData() {
        const now = new Date();
        const day = String(now.getDate()).padStart(2, '0');
        const month = String(now.getMonth() + 1).padStart(2, '0'); // Janeiro é 0!
            const hours = String(now.getHours()).padStart(2, '0');
        const minutes = String(now.getMinutes()).padStart(2, '0');
        return `${day}/${month} ${hours}:${minutes}`;
    }
}
try {
    GerenciadorSala.adicionaSalaMenu("sala__1")
    /*const eventoSSE = {}//new EventSource(`/sse/${iduser}`);
    eventoSSE.close()
    eventoSSE.onerror = function(event) {
        console.error("Erro no SSE: ", event);
    };
    eventoSSE.addEventListener('entrou-chat', e => {
        console.log("TIPO",e.type)
        let {sala, remetente} = JSON.parse(e.data).conteudo
        const room = new Sala(sala)
        room.salva()
        const usuario = Utils.getUsuario(remetente)
        if(!salas[room.id]) {
            GerenciadorSala.adicionaSalaMenu(room.id)
        }
        usuario.salva()
        return
    })
    eventoSSE.addEventListener('chat-nova-mensagem', e => {
        e.data = JSON.parse(e.data)
        let {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
        const room = new Sala(sala.nome);
        const usuario = Utils.getUsuario(remetente)
        const message = new Mensagem({nomeUsuario: usuario.nome, mensagem, idsala: room.nome})
        message.salva()
        room.adicionaMensagemEmContainer(message)
        return
    })
    eventoSSE.addEventListener('chat', e => {
        console.log(e.type)
        console.log(e.data)
        return
    })*/
} catch (error) {
    console.log(error)
}
