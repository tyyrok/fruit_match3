package routes

type Message struct {
	Type string `json:"type"`
	Data map[string]any `json:"data"`
}


type GameBoard struct {
	Cells [8][8]int `json:"cells"`
}
