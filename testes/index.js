const BASE_URL = "http://localhost:3000"
const NOME_SALA_TESTE = "salateste"

const usuarios = [
    //"kadu",
    "nathy",
    "Beto",
    "walter",
    "usuario1",
    "usuario2",
    "usuario3",
    "usuario4",
    "usuario5",
    /*"usu6",
    "usu7",
    "usu8",
    "usu9",
    "usu10",
    "usu11",
    "usu12",
    "usu13",
    "usu14",
    "usu15",
    "usu16",
    "usu17",
    "usu18",
    "usu19",
    "usu20",
    "usu21",
    "usu22",
    "usu23",
    "usu24",
    "usu25",
    "usu26",
    "usu27",
    "usu28",
    "usu29",
    "usu30",
    "usu31",
    "usu32",
    "usu33",
    "usu34",
    "usu35",
    "usu36",
    "usu37",
    "usu38",
    "usu39",
    "usu40",
    "usu41",
    "usu42",
    "usu43",
    "usu44",
    "usu45",
    "usu46",
    "usu47",
    "usu48",
    "usu49",
    "usu50",
    "usu51",
    "usu52",
    "usu53",
    "usu54",
    "usu55",
    "usu56",
    "usu57",
    "usu58",
    "usu59",
    "usu60",
    "usu61",
    "usu62",
    "usu63",
    "usu64",
    "usu65",
    "usu66",
    "usu67",
    "usu68",
    "usu69",
    "usu70",
    "usu71",
    "usu72",
    "usu73",
    "usu74",
    "usu75",
    "usu76",
    "usu77",
    "usu78",
    "usu79",
    "usu80",
    "usu81",
    "usu82",
    "usu83",
    "usu84",
    "usu85",
    "usu86",
    "usu87",
    "usu88",
    "usu89",
    "usu90",
    "usu91",
    "usu92",
    "usu93",
    "usu94",
    "usu95",
    "usu96",
    "usu97",
    "usu98",
    "usu99",
    "usu100",
    "usu101",
    "usu102",
    "usu103",
    "usu104",
    "usu105",
    "usu106",
    "usu107",
    "usu108",
    "usu109",
    "usu110",
    "usu111",
    "usu112",
    "usu113",
    "usu114",
    "usu115",
    "usu116",
    "usu117",
    "usu118",
    "usu119",
    "usu120",
    "usu121",
    "usu122",
    "usu123",
    "usu124",
    "usu125",
    "usu126",
    "usu127",
    "usu128",
    "usu129",
    "usu130",
    "usu131",
    "usu132",
    "usu133",
    "usu134",
    "usu135",
    "usu136",
    "usu137",
    "usu138",
    "usu139",
    "usu140",
    "usu141",
    "usu142",
    "usu143",
    "usu144",
    "usu145",
    "usu146",
    "usu147",
    "usu148",
    "usu149",
    "usu150",
    "usu151",
    "usu152",
    "usu153",
    "usu154",
    "usu155",
    "usu156",
    "usu157",
    "usu158",
    "usu159",
    "usu160",
    "usu161",
    "usu162",
    "usu163",
    "usu164",
    "usu165",
    "usu166",
    "usu167",
    "usu168",
    "usu169",
    "usu170",
    "usu171",
    "usu172",
    "usu173",
    "usu174",
    "usu175",
    "usu176",
    "usu177",
    "usu178",
    "usu179",
    "usu180",
    "usu181",
    "usu182",
    "usu183",
    "usu184",
    "usu185",
    "usu186",
    "usu187",
    "usu188",
    "usu189",
    "usu190",
    "usu191",
    "usu192",
    "usu193",
    "usu194",
    "usu195",
    "usu196"*/
];
const criarClientes = async () => {
    // TENHO QUE ESPERAR POIS O GOLANG RECLMA DE concurrent map writes
    for await (let apelido of usuarios) {
        let url = `${BASE_URL}/sse/`+apelido
        console.log("abrindo sse para ", url)
        fetch(url)
        console.log("esperando... CUNCURRENT MAP WRITES")
      //  await delay(35)
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
const testarEnvioDeMensagens = (limite = 10) => {
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
 //   await delay(3000)
    console.log("\n")
    testarEnvioDeMensagens(1000)
    console.log("Teste passou")
}
main()
