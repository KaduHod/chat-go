const eventoSSE = new EventSource("/event-stream")
eventoSSE.onmessage = function (event) {
    console.log({data: JSON.parse(event.data), "data-2": new Date().getTime()});
    document.getElementById("data").innerText = event.data;
};

eventoSSE.onerror = e => {
    console.error("Error occurred:", e);
};

eventoSSE.onopen = e => {
    console.log("Connection opened:", e);
};
