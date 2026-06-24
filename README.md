# GopherToe — Jogo da Velha via TCP (Go)

Eu desenvolvi este projeto para demonstrar, de forma prática e simples, como construir uma aplicação cliente-servidor em Go que permita a duas pessoas jogarem Tic-Tac-Toe (Jogo da Velha) usando sockets TCP.

## **Visão geral**

O propósito é oferecer um exemplo didático de comunicação em rede, serialização em JSON e uso básico de concorrência em Go. A aplicação é orientada por linha de comando: um servidor aceita duas conexões e coordena a partida entre os jogadores `X` e `O`.

## **Funcionalidades**

- Partida para dois jogadores via TCP.
- Sincronização do estado do jogo (struct `GameState`) em JSON entre servidor e clientes.
- Validação de jogadas tanto no cliente quanto no servidor.
- Encerramento automático quando há vitória ou empate.

## **Arquitetura**

O projeto adota uma arquitetura cliente‑servidor simples e explícita:

- Servidor (`server/main.go`): escuta uma porta TCP, aceita duas conexões sequenciais, mantém o estado central do jogo, envia atualizações de estado a ambos os clientes e aplica jogadas recebidas através de um canal (`chan types.Play`).
- Cliente (`client/main.go`): conecta-se ao servidor, recebe seu identificador (`X` ou `O`), exibe o tabuleiro no terminal, lê a entrada do usuário, valida localmente e envia a jogada ao servidor.

Cada conexão é tratada por uma goroutine (`handleConnection`) que lê linhas JSON e publica `Play` no canal central. Usei `context.Context` para sinalizar encerramento das goroutines quando a partida termina.

## **Comunicação TCP (detalhes técnicos)**

- Mensagens terminadas por newline (`\n`): o servidor envia o `GameState` serializado em JSON seguido de `\n`; o cliente lê com `ReadString('\n')` e faz `json.Unmarshal`.
- Jogadas do cliente: o cliente envia um JSON com `Play` (campos `position` e `player`) seguido de `\n`.
- Handshake inicial: após os dois `Accept()` o servidor envia mensagens de texto simples como "The match started" e o caractere do jogador (`X` ou `O`).

Observação: optei por uma abordagem simples (linha-por-linha) em vez de um protocolo com length-prefix ou cabeçalho — é suficiente para este exemplo, mas tem limitações em cenários mais complexos.

## **Responsabilidades (resumidas)**

- Servidor: fonte única da verdade do `GameState`, aplicação de jogadas via `logic.MakeMove`, verificação de vitória via `logic.CheckVictory`, envio de atualizações a ambos os clientes.
- Cliente: interface com o jogador, validação preliminar de entrada (`logic.CheckMove`), serialização e envio de jogadas.

## **Estado da partida**

O jogo é representado pela struct `GameState` (arquivo `types/structs.go`) com os campos:

- `Matrix` ([3][3]string): o tabuleiro; células vazias contêm uma string com espaço (`" "`).
- `Status` (string): "playing" ou "stopped".
- `Turn` (int): contador de turnos.
- `Player` (string): jogador atual esperado (`X` ou `O`).
- `Winner` (string): `X`, `O` ou " " quando não há vencedor.

O servidor atualiza esses campos e envia cópias para os clientes a cada iteração do loop principal.

## **Tratamento de erros e limitações observadas**

- Validações críticas (por exemplo, porta inválida) encerram a aplicação com `log.Fatal`.
- Em `handleConnection`, erros de leitura (p.ex. desconexão) fazem a goroutine encerrar; erros são impressos com `fmt.Print`.
- Jogadas inválidas são detectadas por `logic.MakeMove` no servidor; o servidor registra "Invalid move" e reenvia o mesmo `GameState` para que o jogador tente novamente.

Limitações relevantes que observei no código:

- Não existe um protocolo de erro estruturado entre servidor e cliente — o feedback de erro depende do reenvio do estado e logs no servidor.
- O servidor aceita exatamente duas conexões por execução; não há suporte a múltiplas partidas concorrentes ou fila de espera.
- O cliente só aceita flag `-port` (não há flag `-host`), logo assume `localhost` como destino.

## **Estrutura do repositório**

- `client/main.go` — cliente CLI que conecta ao servidor, exibe o tabuleiro e envia jogadas.
- `server/main.go` — servidor TCP que orquestra a partida.
- `server/handleConnection.go` — leitura por conexão, desserialização de `Play` e `sendGameState`.
- `logic/main.go` — regras do jogo: impressão do tabuleiro, validação de movimento e verificação de vitória.
- `types/structs.go` — definições de `GameState` e `Play` usadas na (de)seriação JSON.
- `go.mod` — módulo `tictactoe`, declaração de versão do Go.

## **Tecnologias**

| Tecnologia | Uso |
| --- | --- |
| Go (std) | Rede TCP (`net`), JSON (`encoding/json`), concorrência (goroutines, channels), leitura (`bufio`) |

Não utilizei dependências externas neste projeto.

## **Como executar (passo a passo)**

1. Compilar:

```bash
cd /home/joao_vitor/Documentos/programacao/GopherToe
go build ./...
```

2. Iniciar o servidor (porta 8080 por padrão):

```bash
./server -port 8080
```

3. Em dois terminais separados, iniciar os clientes:

```bash
./client -port 8080
```

Quando ambos os clientes estiverem conectados, o servidor envia "The match started" e o identificador do jogador. O cliente ativo verá o tabuleiro e deverá inserir uma posição entre 1 e 9.

Também é possível executar sem compilar usando `go run`:

```bash
go run ./server -port 8080
go run ./client -port 8080
```

## **Fluxo da aplicação (resumido)**

1. `server` aceita duas conexões (player_one e player_two).
2. O servidor envia mensagens iniciais e os identificadores `X`/`O`.
3. Cada conexão tem uma goroutine que lê `Play` em JSON e publica no canal central.
4. O loop principal do servidor envia `GameState` a ambos, aguarda `play := <-channel`, tenta aplicar a jogada, checa vitória e incrementa `Turn`.
5. Quando há vencedor ou empate (`Turn` >= 9), o servidor define `Status = "stopped"`, envia o estado final e cancela o `context`.

## **Exemplo rápido**

Terminal A (servidor):

```bash
go run ./server -port 8080
# Listening on port 8080...
# Server ready on port 8080
```

Terminal B (cliente X):

```bash
go run ./client -port 8080
# Connecting to port 8080...
# connected to port 8080
# waiting for player two
# The match started
# X
# You are the player X
```

Terminal C (cliente O):

```bash
go run ./client -port 8080
# ...recebe O e começa a receber estados em JSON
```

Os clientes trocam jogadas enviando JSON como `{"position":5,"player":"X"}` seguido de newline.

## **Desafios que enfrentei (e que talvez você enfrente)**

- Garantir consistência do estado com concorrência: usei um canal central e uma goroutine por conexão para simplificar a sincronização.
- Tratamento de desconexões e encerramento limpo: apliquei `defer conn.Close()` e `context` para sinalizar término.
- Escolha do formato de mensagens: optei por JSON linha-terminada por simplicidade, ciente das limitações em produção.

## **Aprendizados**

- Como usar `net` e `bufio` para comunicação TCP simples em Go.
- Coordenação entre goroutines via `chan` para receber jogadas de múltiplas conexões.
- Organização de um pequeno projeto Go com responsabilidades separadas (`server`, `client`, `logic`, `types`).

## **Melhorias futuras**

- Implementar mensagens de controle/erro estruturadas (ex.: envelope com tipo e payload).
- Suportar múltiplas partidas simultâneas e emparelhamento de jogadores.
- Adicionar testes unitários para o pacote `logic`.
- Permitir configuração de host no cliente (flag `-host`).

## **Autor**

Sou o autor deste projeto. Você pode me encontrar no GitHub: [Github](https://github.com/JoaoVitor197843)

Repositórios em destaque no meu perfil:

- [GopherToe](https://github.com/JoaoVitor197843/GopherToe) — Implementação em Go deste projeto.
- [AlgorithmsPython](https://github.com/JoaoVitor197843/AlgorithmsPython) — Testes e implementações de algoritmos em Python.
- [Campo_Minado](https://github.com/JoaoVitor197843/Campo_Minado) — Campo minado implementado em Python puro.
- [Bot-Discord-Ninja-RPG](https://github.com/JoaoVitor197843/Bot-Discord-Ninja-RPG) — Bot para Discord escrito em Python.
- [whyle-finance](https://github.com/JoaoVitor197843/whyle-finance) — Plataforma de gerenciamento financeiro (TypeScript).
- [joao-portfolio](https://github.com/JoaoVitor197843/joao-portfolio) — Portfolio minimalista (Next.js, Tailwind CSS).
