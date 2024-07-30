const iduser = document.getElementById("idusuario").value;
const salaitemConstructor = nomeSala => `<div id="sala___${nomeSala}" class="item-card rounded">${nomeSala}</div>`
const salas = []
const salasMensagens = {}
const salaItemContainer = document.getElementById("chat-salas")
function formatarData() {
    const now = new Date();
    const day = String(now.getDate()).padStart(2, '0');
    const month = String(now.getMonth() + 1).padStart(2, '0'); // Janeiro é 0!
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    return `${day}/${month} ${hours}:${minutes}`;
}
const chatContainer = document.getElementById("chat-aberto-mensagens")
const montaMensagem =  (msg, remetente) => {
    let classeAlinhamento = remetente == iduser ? "self-end" : "self-start"
	let cor = msg.remetente == "Arnilloy" ? "#46295A" : "#800000";
    return `<div class="${classeAlinhamento} flex shadow-md flex-col rounded-md p-2 max-w-fit break-words h-fit bg-slate-200 my-1">
                <div class="message-content">
                    ${msg}
                </div>
                <div class="flex justify-end mb-2">
                    <span class="mr-1 font-bold" style="color:${cor}">${remetente}</span>
                    <span class="text-slate-600">${formatarData()}</span>
                </div>
            </div>`
}
const eventoSSE = new EventSource(`/sse/${iduser}`);
eventoSSE.onerror = function(event) {
    console.error("Erro no SSE: ", event);
};
eventoSSE.addEventListener('entrou-chat', e => {
    console.log("TIPO",e.type)
    let json = JSON.parse(e.data)
    if (!salas.includes(json.conteudo.sala)) {
        salaItemContainer.innerHTML += salaitemConstructor(json.conteudo.sala)
    }
    return
})
eventoSSE.addEventListener('chat-nova-mensagem', e => {
    e.data = JSON.parse(e.data)
    const {sala, remetente, mensagem} = JSON.parse(e.data).conteudo
    console.log(sala, remetente, mensagem)
    div = montaMensagem(mensagem, remetente)
    chatContainer.innerHTML += div
    salasMensagens[sala].push(div)
    return
})
eventoSSE.addEventListener('chat', e => {
    console.log(e.type)
    console.log(e.data)
    return
})
