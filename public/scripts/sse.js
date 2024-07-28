document.addEventListener("DOMContentLoaded", function() {
    const iduser = document.getElementById("idusuario").value;
    const eventoSSE = new EventSource(`/sse/${iduser}`);

    eventoSSE.onmessage = function(event) {
        console.log("Mensagem recebida: ", event.data);

        try {
            const dadosEvento = JSON.parse(event.data);
            const {tipo, conteudo} = dadosEvento;
            console.log("Dados do evento", {tipo, conteudo});
            document.getElementById("data").innerText = JSON.stringify(dadosEvento, null, 2);
        } catch (e) {
            console.error("Erro ao processar dados do evento", e);
        }
    };

    eventoSSE.addEventListener("message", function(event) {
        console.log("Evento 'message' recebido: ", event.data);
    });

    eventoSSE.onerror = function(event) {
        console.error("Erro no SSE: ", event);
    };
});
