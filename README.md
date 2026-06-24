# GopherToe

Jogo da velha multiplayer jogado por dois clientes conectados a um servidor via TCP. O nome é um trocadilho entre o **Gopher** (mascote do Go) e *Tic-Tac-**Toe***.

Escrevi o projeto para praticar, na prática, três pilares de backend em Go: comunicação em rede com sockets TCP, concorrência com goroutines e channels, e serialização de estado em JSON — sem nenhuma dependência externa, só a biblioteca padrão.

## Como funciona

O servidor escuta uma porta TCP e aceita exatamente duas conexões. O primeiro jogador recebe `X` e fica aguardando; quando o segundo entra (recebe `O`), a partida começa. A partir daí o servidor é a **fonte única da verdade**: ele mantém o estado do jogo, recebe as jogadas, valida, atualiza o tabuleiro e devolve o estado atualizado para os dois clientes a cada turno, até dar vitória ou empate.

Cada conexão é tratada por uma goroutine que lê as jogadas (JSON) e as publica em um channel central. O loop principal do servidor consome esse channel, aplica a jogada do jogador da vez e sincroniza todo mundo. Isso mantém a lógica de turnos simples e evita condições de corrida sobre o estado compartilhado.

## Stack

Apenas a biblioteca padrão do Go:

- `net` — sockets TCP (servidor e cliente)
- `encoding/json` — serialização do estado e das jogadas
- goroutines e `chan` — uma goroutine por conexão, channel central para coordenar as jogadas
- `context` — sinalização de encerramento das goroutines no fim da partida
- `bufio` — leitura das mensagens linha a linha

## Estrutura

```
GopherToe/
├── server/
│   ├── main.go              # orquestra a partida: aceita conexões, mantém o estado, loop principal
│   └── handleConnection.go  # uma goroutine por jogador: lê jogadas e publica no channel
├── client/
│   └── main.go              # conecta, recebe o estado, lê a jogada do usuário e envia
├── logic/
│   └── main.go              # regras puras: impressão do tabuleiro, validação e vitória
├── types/
│   └── structs.go           # GameState e Play (compartilhados entre cliente e servidor)
└── go.mod                   # módulo tictactoe
```

Separei `types` e `logic` em pacotes próprios para que cliente e servidor compartilhem exatamente as mesmas structs e regras, sem duplicação.

## Protocolo

A comunicação é por **mensagens JSON terminadas em newline** (`\n`). Cada lado lê com `ReadString('\n')` e faz `json.Unmarshal`.

- **Servidor → cliente:** o `GameState` completo a cada turno (tabuleiro, status, vez, vencedor).
- **Cliente → servidor:** uma jogada — `{"position": 5, "player": "X"}`.

Optei por essa abordagem simples (uma mensagem por linha) em vez de um protocolo com cabeçalho ou *length-prefix*. É suficiente para o escopo do jogo, mas tenho consciência de que não cobre cenários mais complexos (mensagens fragmentadas, payloads grandes).

## Como rodar

Pré-requisito: Go instalado.

```bash
git clone https://github.com/JoaoVitor197843/GopherToe
cd GopherToe

# Terminal 1 — servidor (porta 8080 por padrão)
go run ./server -port 8080

# Terminal 2 — jogador X
go run ./client -port 8080

# Terminal 3 — jogador O
go run ./client -port 8080
```

Quando os dois clientes conectam, a partida começa. O jogador da vez escolhe uma posição de 1 a 9 no tabuleiro.

## Decisões e limitações

Algumas escolhas conscientes, com seus trade-offs:

- **Channel central em vez de locks manuais** para o estado compartilhado — mais simples de raciocinar e suficiente para dois jogadores.
- **Sem protocolo de erro estruturado:** quando uma jogada é inválida, o servidor apenas reenvia o mesmo estado para o jogador tentar de novo, em vez de mandar uma mensagem de erro dedicada.
- **Uma partida por execução:** o servidor atende exatamente dois jogadores; não há suporte a partidas simultâneas ou fila de espera.
- **Cliente assume `localhost`:** existe apenas a flag `-port`, não `-host`.

## Próximos passos

- Envelopar as mensagens em um formato com tipo + payload (separar estado, erro e controle).
- Suportar múltiplas partidas simultâneas com pareamento de jogadores.
- Testes unitários para o pacote `logic`.
- Flag `-host` no cliente para jogar fora do `localhost`.

## Autor

João Vitor — [github.com/JoaoVitor197843](https://github.com/JoaoVitor197843)
