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

func (g greg) Do(fields [9]pkg.FieldType, round int) (field uint8, action pkg.Action) {
	field = uint8(r.Intn(9))
	action = pkg.ResolveAction(fields[field])

	return
}

var Greg greg
