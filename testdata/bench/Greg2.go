package main

import (
	"github.com/gregmus2/ants-pkg"
	"math/rand"
)

type greg2 string

var r2 *rand.Rand

func init() {
	r2 = rand.New(rand.NewSource(555))
}

func (g greg2) Start(anthill pkg.Pos, birth pkg.Pos) {

}

func (g greg2) Do(fields [5][5]pkg.FieldType, round int) (target *pkg.Pos, action pkg.Action) {
	target = &pkg.Pos{X: r2.Intn(3) - 1, Y: r2.Intn(3) - 1}
	action = pkg.ResolveAction(fields[target.X + 3][target.Y + 3])

	return
}

var Greg2 greg2
