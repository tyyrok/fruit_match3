package routes

type Message struct {
	Type string `json:"type"`
	Data map[string]any `json:"data"`
}

type GameBoard struct {
	Cells [8][8]int `json:"cells"`
	Scores int `json:"scores"`
}

type Turn struct {
	FromRow, ToRow, FromCol, ToCol int
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *Point) equal(d *Point) bool {
	if p.X == d.X && p.Y == d.Y {
		return true
	} else {
		return false
	}
}

type Combination struct {
	Points []Point `json:"points"`
}

func (c *Combination) equal(d *Combination) bool {
	if len(c.Points) != len(d.Points) {
		return false
	}
	for i := 0; i < len(c.Points); i++ {
		if !c.Points[i].equal(&d.Points[i]) {
			return false
		}
	}
	return true
}

func (c *Combination) getLenght() int {
	return len(c.Points)
}

func (c *Combination) checkIntersection(d *Combination) bool {
	for _, val_c := range c.Points {
		for _, val_d := range d.Points {
			if val_c.equal(&val_d) {
				return true
			}
		}
	}
	return false
}
