package logic

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Position struct {
	Line   int
	Column int
}

func PrintMatrix(matrix [3][3]string) {
	position_number := 1
	for row := range matrix {
		for column := range matrix[row] {
			if matrix[row][column] != " " {
				fmt.Printf("%s", matrix[row][column])
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
func CheckMove(matrix *[3][3]string, row, column int, player string) error {
	if row < 0 || row >= len(matrix) {
		return fmt.Errorf("row %d is out of bounds", row)
	}
	if column < 0 || column >= len(matrix[row]) {
		return fmt.Errorf("column %d is out of bounds", column)
	}
	if player != "X" && player != "O" {
		return fmt.Errorf("player %s not accepted", player)
	}
	if matrix[row][column] != " " {
		return fmt.Errorf("a move has already been made from this position")
	}
	return nil
}
func MakeMove(matrix *[3][3]string, row, column int, player string) error {
	if row < 0 || row >= len(matrix) {
		return fmt.Errorf("row %d is out of bounds", row)
	}
	if column < 0 || column >= len(matrix[row]) {
		return fmt.Errorf("column %d is out of bounds", column)
	}
	if player != "X" && player != "O" {
		return fmt.Errorf("player %s not accepted", player)
	}
	if matrix[row][column] != " " {
		return fmt.Errorf("a move has already been made from this position")
	}
	matrix[row][column] = player
	return nil
}
func CheckVictory(matrix [3][3]string) string {
	victoryConditions := [8][3]Position{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},

		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},

		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}

	for _, line := range victoryConditions {
		a := matrix[line[0].Line][line[0].Column]
		b := matrix[line[1].Line][line[1].Column]
		c := matrix[line[2].Line][line[2].Column]
		if a != " " && a == b && b == c {
			return a
		}
	}

	return " "
}

func ReadNumber(reader *bufio.Reader, prompt string) (int, error) {
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
