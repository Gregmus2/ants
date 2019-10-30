package main

import (
	"github.com/gregmus2/ants-pkg"
	"math/rand"
	"time"
)

type greg string

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {

}

func (g greg) Start(anthill pkg.Pos, birth pkg.Pos) {

}

func (g greg) Do(fields [5][5]pkg.FieldType, round int) (target pkg.Pos, action pkg.Action) {
	target = pkg.Pos{r.Intn(5), r.Intn(5)}
	action = pkg.ResolveAction(fields[target.X()][target.Y()])

	return
}

var Greg greg
