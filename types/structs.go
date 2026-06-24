package types

type GameState struct {
	Matrix [3][3]string `json:"matrix"`
	Status string       `json:"status"`
	Turn   int          `json:"turn"`
	Player string       `json:"player"`
	Winner string       `json:"winner"`
}
type Play struct {
	Position int    `json:"position"`
	Player   string `json:"player"`
}
