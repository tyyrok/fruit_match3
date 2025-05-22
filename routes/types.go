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


func (b *GameBoard) updateBoadByTurn(t *Turn) {
	b.Cells[t.FromRow][t.FromCol], b.Cells[t.ToRow][t.ToCol] = b.Cells[t.ToRow][t.ToCol], b.Cells[t.FromRow][t.FromCol]
}

func (b *GameBoard) updateBoard(combs *[]Combination, newElems []int) {
	for _, comb := range *combs {
		if comb.isHorizontal() {
			for _, point := range comb.Points {
				i := point.Y
				for i > 0 {
					b.Cells[i][point.X] = b.Cells[i-1][point.X]
					i -= 1
				}
				b.Cells[0][point.X] = newElems[0]
				newElems = newElems[1:]
			}
		} else {
			for _, point := range comb.Points {
				b.Cells[point.Y][point.X] = newElems[0]
				newElems = newElems[1:]
			}
		}
		
	}
}

func (c *Combination) isHorizontal() bool {
	return c.Points[0].Y == c.Points[1].Y
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
