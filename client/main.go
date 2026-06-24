package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"tictactoe/logic"
	"tictactoe/types"
)

func main() {
	port := flag.Int("port", 8080, "Server connection port")
	flag.Parse()
	if *port < 1 || *port > 65535 {
		log.Fatal("Invalid port")
	}
	fmt.Printf("Connecting to port %d...\n", *port)
	selected_port := fmt.Sprintf("localhost:%d", *port)
	server, err := net.Dial("tcp", selected_port)
	if err != nil {
		log.Fatalf("can't connect to port %d\n%s", *port, err)
	}
	defer server.Close()
	fmt.Printf("connected to port %d\n", *port)
	server_reader := bufio.NewReader(server)
	for messages := 0; messages < 2; messages++ {
		msg, err := server_reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		fmt.Print(msg)
		if msg == "The match started\n" {
			break
		}
	}
	player, err := server_reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	player = strings.TrimSpace(player)
	fmt.Printf("You are the player %s\n", player)
	var gameState types.GameState
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := server_reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		json.Unmarshal([]byte(line), &gameState)
		if gameState.Status != "playing" {
			switch gameState.Winner {
			case player:
				fmt.Println("You win!")
			case " ":
				fmt.Println("That's a draw!")
			default:
				fmt.Println("You lose!")
			}
			break
		}
		logic.PrintMatrix(gameState.Matrix)
		if gameState.Player == player {
			fmt.Println("Your turn")
			var number int
			for {
				number, err = logic.ReadNumber(reader, "enter a position between 1 and 9: ")
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = logic.CheckMove(&gameState.Matrix, (number-1)/3, (number-1)%3, player)
				if err != nil {
					fmt.Println(err)
					continue
				}
				break
			}
			play := types.Play{
				Position: number,
				Player:   player,
			}
			data, err := json.Marshal(play)
			if err != nil {
				fmt.Print(err)
				return
			}
			server.Write(append(data, '\n'))
		} else {
			fmt.Println("your opponent's turn")
		}
	}
	server.Close()
}
