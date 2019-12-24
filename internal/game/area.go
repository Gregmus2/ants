package game

import (
	pkg "github.com/gregmus2/ants-pkg"
)

type Area struct {
	w, h   int
	matrix [][]*Object
}

func NewArea(w, h int) *Area {
	matrix := make([][]*Object, w)
	for x := 0; x < w; x++ {
		matrix[x] = make([]*Object, h)
		for y := 0; y < h; y++ {
			// edges
			if x == 0 || x == w-1 || y == 0 || y == h-1 {
				matrix[x][y] = CreateWall()
			} else {
				matrix[x][y] = CreateEmptyObject()
			}
		}
	}

	return &Area{
		w:      w,
		h:      h,
		matrix: matrix,
	}
}

func (a *Area) ToColorSlice() [][]string {
	colorSlice := make([][]string, len(a.matrix))
	for x := 0; x < a.w; x++ {
		colorSlice[x] = make([]string, a.h)
		for y := 0; y < a.h; y++ {
			colorSlice[x][y] = a.matrix[x][y].Color
		}
	}

	return colorSlice
}

func (a *Area) VisibleArea(ant *Ant) [5][5]pkg.FieldType {
	var fieldTypes [5][5]pkg.FieldType
	for dY := 0; dY < 5; dY++ {
		for dX := 0; dX < 5; dX++ {
			x := ant.Pos.X + dX - 2
			y := ant.Pos.Y + dY - 2
			if x < 0 || y < 0 || x > a.w-1 || y > a.h-1 {
				fieldTypes[dX][dY] = pkg.NoField
				continue
			}

			fieldTypes[dX][dY] = a.matrix[x][y].FieldTypeForUser(ant)
		}
	}

	return fieldTypes
}

func (a *Area) CalcAtkPower(target *Ant, attacker *Ant) int {
	power := 0
	for y := target.Pos.Y - 1; y <= target.Pos.Y+1; y++ {
		for x := target.Pos.X - 1; x <= target.Pos.X+1; x++ {
			if a.matrix[x][y].Type != pkg.AntField {
				continue
			}

			switch a.matrix[x][y].Ant.User {
			case target.User:
				power--
			case attacker.User:
				power++
			}
		}
	}

	return power
}

func (a *Area) ByPos(pos *pkg.Pos) *Object {
	return a.matrix[pos.X][pos.Y]
}
