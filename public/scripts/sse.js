const iduser = document.getElementById("idusuario").value;
const salaItemContainer = document.getElementById("chat-salas")
const chatContainer = document.getElementById("chat-aberto-mensagens")
const listaMensagens = [];
const listaUsuarios = [];
const listaSalas = [];
const listaSalasUsuarios = [];
const listaSalasUsuariosMensagens = [];
class Utils {
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
class Usuario {
    constructor(nome, id = null, cor = null) {
        this.id = id ?? "usuario__"+Utils.generateUniqueId()
        this.nome = nome
        this.cor = cor ?? Utils.getRandomHexColor()
    }
    ehUsuarioLogado() {
        return this.nome === document.getElementById("idusuario").value;
    }
}
class Mensagem {
    constructor(mensagem, idremetente){
        this.mensagem = mensagem
        this.id = "mensagem__"+Utils.generateUniqueId()
        this.data = Utils.formataData();
        this.idremetente = idremetente
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
}
class Sala {
    constructor(nomeSala){
        this.nome = nomeSala
        this.id = `sala__${this.nome}`
    }
    montaHtml(){
        this.html = `<div
            id="${this.id}"
            class="item-card rounded"
        >${this.nome}</div>`
        return this.html
    }
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
}
try {
    const eventoSSE = {}//new EventSource(`/sse/${iduser}`);
    eventoSSE.onerror = function(event) {
        console.error("Erro no SSE: ", event);
    };
    eventoSSE.addEventListener('entrou-chat', e => {
        const {sala, remetente} = JSON.parse(e.data).conteudo
        let room = listaSalas.find(room => room.nome === sala)
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
    })
    eventoSSE.addEventListener('chat-nova-mensagem', e => {
        e.data = JSON.parse(e.data)
        let {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
        const user = listaUsuarios.find(user => user.nome === remetente)
        if(!user) {
            throw { error: "Usuario não entrou na sala" }
        }
        const room = listaSalas.find(room => room.nome === sala)
        if(room) {
            throw { error: "Sala não encontrada"}
        }
        const message = new Mensagem(mensagem, user.id)
        listaMensagens.push(message)
        const salaUsuarioMensagem = new SalaUsuarioMensagem(room.id, user.id, message.id)
        listaSalasUsuariosMensagens.push(salaUsuarioMensagem)
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
