document.addEventListener("DOMContentLoaded", function() {
    const iduser = document.getElementById("idusuario").value;
    const eventoSSE = new EventSource(`/sse/${iduser}`);
    eventoSSE.onerror = function(event) {
        console.error("Erro no SSE: ", event);
    };
    eventoSSE.addEventListener('ping', e => {
    })
    eventoSSE.addEventListener('entrou-chat', e => {
        console.log(e.type)
        console.log(e.data)
    })
    eventoSSE.addEventListener('chat-nova-mensagem', e => {
        console.log(e.type)
        console.log(e.data)
    })
    eventoSSE.addEventListener('chat', e => {
        console.log(e.type)
        console.log(e.data)
    })
    console.log(eventoSSE)
});
