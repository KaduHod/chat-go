**WebChat em Go**
=====================

Um aplicativo de conversa em tempo real escrito em Golang, utilizando SSE (Server-Sent Events) e canais em Go para gerar o
motor de mensagens. Com esse aplicativo, é possível criar um ambiente de conversa online com amigos ou familiares.

### Características:

* Motor de mensagens personalizado
* Utilização de SSE (Server-Sent Events) para atualizar a página em tempo real
* Canais em Go para gerar o motor de mensagens
* Fácil de usar: basta digitar `docker compose up -d` para iniciar o programa e `go build -buildvcs=false . && /chat` dentro
do container

### Como usar:

1. Inicie o aplicativo com `docker compose up -d`
2. Entre no container com `docker exec -it <container_name> sh`
3. Digite `go build -buildvcs=false . && /chat` para iniciar o programa
4. Comece a chatear!

**1. `/chat/:apelidousuario/salas`**

* **Method:** GET
* **Path:** `/chat/{apelidousuario}/salas`
* **Description:** Retorna as salas de um usuário especificado.
* **Parameters:**
	+ `{apelidousuario}` (string): O apelido do usuário para buscar.
* **Response:**
	+ 200 OK: Uma lista de salas do usuário, caso o usuário seja encontrado. Caso contrário, uma mensagem de erro é retornada.
	+ Exemplo de resposta:
```json
[
    {
        "nome": "Sala de Bicicleta",
        "id": 1
    },
    {
        "nome": "Sala de Futebol",
        "id": 2
    }
]
```
* **Error Response:**
	+ 404 Not Found: O usuário não foi encontrado.
```json
{
    "error": "Usuário não encontrado"
}
```

**2. `/chat/sala/:nomesala/usuarios`**

* **Method:** GET
* **Path:** `/chat/sala/{nomesala}/usuarios`
* **Description:** Retorna os usuários de uma sala especificada.
* **Parameters:**
	+ `{nomesala}` (string): O nome da sala para buscar.
* **Response:**
	+ 200 OK: Uma lista de usuários da sala, caso a sala seja encontrada. Caso contrário, uma mensagem de erro é retornada.
	+ Exemplo de resposta:
```json
[
    {
        "apelido": "johnDoe",
        "id": 1
    },
    {
        "apelido": "janeSmith",
        "id": 2
    }
]
```
* **Error Response:**
	+ 404 Not Found: A sala não foi encontrada.
```json
{
    "error": "Sala não encontrada"
}
```

### Rota: `/chat/sse/:apelidousuario/sala/:nomesala/enviar`

* **Método:** POST
* **Parâmetros:**
	+ `apelidousuario`: O apelido do usuário que está enviando a mensagem.
	+ `nomesala`: O nome da sala para a qual a mensagem está sendo enviada.
	+ `msg`: O conteúdo da mensagem (obrigatório).
* **Descrição:** Esta rota permite ao usuário enviar uma mensagem para uma sala específica. O usuário deve estar conectado à
canal SSE e ter permissão para enviar mensagens na sala.
* **Respostas:**
	+ 201 Created: A mensagem foi enviada com sucesso.
	+ 400 Bad Request: Erro de parâmetro, como falta do conteúdo da mensagem ou usuário não é membro ativo da sala.

### Rota: `/chat/sse/:apelidousuario/sala/:nomesala`

* **Método:** GET
* **Parâmetros:**
	+ `apelidousuario`: O apelido do usuário que está entrando na sala.
	+ `nomesala`: O nome da sala para a qual o usuário está tentando entrar.
* **Descrição:** Esta rota permite ao usuário entrar em uma sala específica. O usuário deve estar conectado à canal SSE e ter
permissão para entrar na sala.
* **Respostas:**
	+ 200 OK: O usuário entrou na sala com sucesso.
	+ 400 Bad Request: Erro de parâmetro, como usuário não está conectado à canal SSE ou não tem permissão para entrar na sala.

A rota `/sse/:apelidousuario` é responsável por estabelecer uma conexão SSE (Server-Sent Events) com o usuário e gerenciar a
transmissão de eventos entre o servidor e o cliente.

Essa rota realiza as seguintes ações:

1. Busca o usuário pelo apelido fornecido no parâmetro `apelidousuario`.
2. Verifica se o usuário já tem uma conexão SSE aberta.
	* Se sim, retorna um erro (400) com a mensagem "usuario já possui sse aberto".
3. Busca o registro do usuário no banco de dados para obter informações como o ID e nome do usuário.
4. Cria um novo canal SSE para o usuário e registra-o no gerenciador de canais SSE.
5. Estabelece a conexão SSE com o cliente, definindo os cabeçalhos necessários (Content-Type, Cache-Control e Connection).
6. Inicia um processo em segundo plano que monitora a conexão SSE do cliente e fecha o canal quando o cliente fecha a
conexão.

**Rota:** `/chat/api/usuario/login`

* **Método:** POST
* **Parâmetros:**
    **Request Body:**
    ```json
    {
      "nome": "<name>",
    }
    ```

* **Descrição:** Esta rota é responsável por lidar com a autenticação de um usuário. Ela recebe os dados de login e verifica se o usuário existe no banco de dados.
* **Respostas:**
	+ 200 OK: O usuário foi autenticado com sucesso.
	+ 404 Not Found: O usuário não existe ou não foi encontrado no banco de dados.
	+ 500 Internal Server Error: Erro ao buscar usuário no banco de dados.
