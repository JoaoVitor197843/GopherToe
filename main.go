package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	Line int
	Column int
}

func printMatrix(matrix [3][3]rune) {
	position_number := 1
	for row := range matrix {
		for column := range matrix[row] {
			if matrix[row][column] != ' ' {
				fmt.Printf("%c", matrix[row][column])
			} else {
			fmt.Printf("%d", position_number)
			}
			if column < 2 {
				fmt.Print(" | ")
			}
			position_number++
		}
		if row < 2 {
			fmt.Println("\n---------")
		}
	}
	fmt.Println()
}

func makeMove(matrix *[3][3]rune, row, column int, player rune) error {
	if row < 0 || row >= len(matrix) {
		return fmt.Errorf("row %d is out of bounds", row)
	}
	if column < 0 || column >= len(matrix[row]) {
		
	return fmt.Errorf("column %d is out of bounds", column)
	}
	if player != 'X' && player != 'O' {
		
		return fmt.Errorf("player %c not accepted", player)
	}
	if matrix[row][column] != ' ' {
		return fmt.Errorf("a move has already been made from this position")
	}
	matrix[row][column] = player
	return nil
}
func checkVictory(matrix [3][3]rune) rune {
	victoryConditions := [8][3]Position{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},

		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},

		{{0, 0},{1, 1},{2, 2}},
		{{0, 2},{1, 1},{2, 0}},
	}

	for _, line := range victoryConditions {
		a := matrix[line[0].Line][line[0].Column]
		b := matrix[line[1].Line][line[1].Column]
		c := matrix[line[2].Line][line[2].Column]
		if a != ' ' && a == b && b == c {
			return a
		}
	}	

	return ' '
}

func readNumber(reader *bufio.Reader, prompt string) (int, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	text = strings.TrimSpace(text)
	number, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid number", text)
	}
	return number, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	matrix := [3][3]rune{
		{' ', ' ', ' '},
		{' ', ' ', ' '},
		{' ', ' ', ' '},
	}

	printMatrix(matrix)
	var victory rune = ' '
	turn := 0
	for victory == ' ' && turn < 9 {
		player := 'X'
		if turn % 2 != 0 {
			player = 'O'
		}
		fmt.Printf("player %c turn\n", player)
		number, err := readNumber(reader, "enter a position between 1 and 9: ")
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = makeMove(&matrix, (number - 1) / 3, (number - 1) % 3, player)
		if err != nil {
			fmt.Println(err)
			continue
		}
		printMatrix(matrix)
		victory = checkVictory(matrix)
		turn++
	}
	if victory != ' ' {
		fmt.Printf("player %c wins\n", victory)
	} else {
		fmt.Println("draw!")
	}	
}