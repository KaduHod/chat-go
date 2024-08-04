const BASE_URL = "http://localhost:3000"
const NOME_SALA_TESTE = "salateste"
let usuarios = [
    "kadu",
    "nathy",
    "Beto",
    "walter"
]
const criarClientes = async () => {
    // TENHO QUE ESPERAR POIS O GOLANG RECLMA DE concurrent map writes
    for await (let apelido of usuarios) {
        let url = `${BASE_URL}/sse/`+apelido
        console.log("abrindo sse para ", url)
        fetch(url)
        console.log("esperando... CUNCURRENT MAP WRITES")
        await delay(35)
    }
}
const entrarEmSalaDeTeste = async () => {
    for await (let apelido of usuarios) {
        console.log({apelido})
        const url = `${BASE_URL}/chat/sse/${apelido}/entrar/${NOME_SALA_TESTE}`
        console.log("Entrando na sala teste", url)
        await fetch(url)

    }
}
const enivarMensagem = (apelido, nomesala, msg, id) => {
    return new Promise(r => {
        const url = `${BASE_URL}/chat/sse/${apelido}/sala/${nomesala}/enviar?msg=${msg}`
        console.log("Enviando mensagem", id)
        fetch(url, {method:"POST"})
        return r(true)
    })
}
const delay = ms => new Promise(resolve => setTimeout(resolve, ms));
const testarEnvioDeMensagens = (limite = 100) => {
    let cont = 0
    while(cont <= limite) {
        console.log("\n")
        usuarios.forEach(apelido => {
            enivarMensagem(apelido, NOME_SALA_TESTE, "Mensagem de teste", cont)
        })
        cont++
    }
}
const main = async () => {
    console.log("Iniciando teste")
    console.log("\n")
    await criarClientes()
    await delay(3000)
    console.log("\n")
    await entrarEmSalaDeTeste()
    await delay(3000)
    console.log("\n")
    testarEnvioDeMensagens(200)
    console.log("Teste passou")
}
main()
