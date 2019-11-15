package game

import (
	pkg "github.com/gregmus2/ants-pkg"
)

type Area [][]*Object

func (a Area) ToColorSlice() [][]string {
	colorSlice := make([][]string, len(a))
	for x := 0; x < len(a); x++ {
		colorSlice[x] = make([]string, len(a[x]))
		for y := 0; y < len(a); y++ {
			colorSlice[x][y] = a[x][y].Color
		}
	}

	return colorSlice
}

func (a Area) VisibleArea(ant *Ant) [5][5]pkg.FieldType {
	fieldTypes := [5][5]pkg.FieldType{}
	for dY := 0; dY < 5; dY++ {
		for dX := 0; dX < 5; dX++ {
			x := ant.Pos.X - dX - 2
			y := ant.Pos.Y - dY - 2
			if x < 0 || y < 0 {
				fieldTypes[dX][dY] = pkg.NoField
				continue
			}

			fieldTypes[dX][dY] = a[x][y].FieldTypeForUser(ant)
		}
	}

	return fieldTypes
}

func (a Area) CalcAtkPower(target *Ant, attacker *Ant) int {
	power := 0
	for y := target.Pos.Y - 1; y <= target.Pos.Y+1; y++ {
		for x := target.Pos.X - 1; x <= target.Pos.X+1; x++ {
			if a[x][y].Type != pkg.AntField {
				continue
			}

			switch a[x][y].Ant.User {
			case target.User:
				power--
			case attacker.User:
				power++
			}
		}
	}

	return power
}

func (a Area) ByPos(pos *pkg.Pos) *Object {
	return a[pos.X][pos.Y]
}
