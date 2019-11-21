package main

import (
	"github.com/gregmus2/ants-pkg"
	"math/rand"
)

type greg string

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(666))
}

func main() {

}

func (g greg) Start(anthill pkg.Pos, birth pkg.Pos) {

}

func (g greg) Do(fields [5][5]pkg.FieldType, round int) (target *pkg.Pos, action pkg.Action) {
	target = &pkg.Pos{r.Intn(3) - 1, r.Intn(3) - 1}
	action = pkg.ResolveAction(fields[target.X + 3][target.Y + 3])

	return
}

var Greg greg
