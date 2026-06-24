package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"tictactoe/logic"
	"tictactoe/types"
)

func main() {
	port := flag.Int("port", 8080, "Server Port")
	flag.Parse()
	if *port < 1 || *port > 65535 {
		log.Fatal("port must be between 1 and 65535")
	}
	fmt.Printf("Listening on port %d...\n", *port)
	selected_port := fmt.Sprintf(":%d", *port)
	listener, err := net.Listen("tcp", selected_port)

	if err != nil {
		log.Fatalf("Can't listen port %d\n%s", *port, err)
	}
	defer listener.Close()
	ctx, close := context.WithCancel(context.Background())
	fmt.Printf("Server ready on port %d\n", *port)
	player_one, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	player_one.Write([]byte("waiting for player two\n"))
	player_two, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	player_one.Write([]byte("The match started\n"))
	player_one.Write([]byte("X\n"))
	player_two.Write([]byte("The match started\n"))
	player_two.Write([]byte("O\n"))
	channel := make(chan types.Play)
	go handleConnection(player_one, channel, "X", ctx)
	go handleConnection(player_two, channel, "O", ctx)
	gameState := types.GameState{
		Matrix: [3][3]string{
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
		},
		Status: "playing",
		Turn:   0,
		Player: "X",
		Winner: " ",
	}
	for gameState.Winner == " " && gameState.Turn < 9 {
		gameState.Player = "X"
		if gameState.Turn%2 != 0 {
			gameState.Player = "O"
		}
		sendGameState(player_one, player_two, gameState)
		play := <-channel
		err := logic.MakeMove(&gameState.Matrix, (play.Position-1)/3, (play.Position-1)%3, play.Player)
		if err != nil {
			fmt.Println("Invalid move")
			continue
		}
		gameState.Winner = logic.CheckVictory(gameState.Matrix)
		gameState.Turn++
	}
	gameState.Status = "stopped"
	sendGameState(player_one, player_two, gameState)
	close()
}
