package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"tictactoe/types"
)

func handleConnection(conn net.Conn, channel chan types.Play, player string, ctx context.Context) {
	reader := bufio.NewReader(conn)
	defer conn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Print(err)
				return
			}
			var play types.Play
			json.Unmarshal([]byte(line), &play)
			if play.Player != player {
				continue
			}
			channel <- play
		}
	}
}

func sendGameState(player_one, player_two net.Conn, gameState types.GameState) {
	data, err := json.Marshal(gameState)
	if err != nil {
		fmt.Print(err)
		return
	}
	player_one.Write(append(data, '\n'))
	player_two.Write(append(data, '\n'))
}
