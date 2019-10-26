package global

import (
	"math"

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

func (a Area) NearestArea(ant *Ant) [9]pkg.FieldType {
	/*  It's fields near ant in that order:
			0 1 2
			3 4 5
	 		6 7 8
	*/
	fieldTypes := [9]pkg.FieldType{}
	i := 0
	for y := ant.Pos.Y() - 1; y <= ant.Pos.Y()+1; y++ {
		for x := ant.Pos.X() - 1; x <= ant.Pos.X()+1; x++ {
			fieldTypes[i] = a[x][y].FieldTypeForUser(ant)
			i++
		}
	}

	return fieldTypes
}

/* add field with format
0 1 2
3 4 5
6 7 8
	to input position
*/
// todo move it to ants-pkg
func (a Area) RelativePosition(pos Pos, field uint8) Pos {
	return Pos{
		pos.X() + uint(math.Mod(float64(field+3), 3)) - 1,
		pos.Y() + uint(math.Floor(float64(field/3))) - 1,
	}
}

func (a Area) ByPos(pos Pos) *Object {
	return a[pos.X()][pos.Y()]
}
