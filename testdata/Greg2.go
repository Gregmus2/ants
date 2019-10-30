package main

import (
	"github.com/gregmus2/ants-pkg"
	"math/rand"
	"time"
)

type greg2 string

var r2 *rand.Rand

func init() {
	r2 = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (g greg2) Start(anthill pkg.Pos, birth pkg.Pos) {

}

func (g greg2) Do(fields [5][5]pkg.FieldType, round int) (target pkg.Pos, action pkg.Action) {
	target = pkg.Pos{r2.Intn(5), r2.Intn(5)}
	action = pkg.ResolveAction(fields[target.X()][target.Y()])

	return
}

var Greg2 greg2
